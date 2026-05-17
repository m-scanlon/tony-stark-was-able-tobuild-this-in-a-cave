package tui

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"skyra-v05/src/reality"
)

var (
	sidebarWidth = 28

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	sidebarStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(1, 1)

	chatStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("8")).
			Padding(0, 1)

	inputStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("12"))

	senderUser = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("12"))

	senderBeing = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("5"))

	dimStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	activeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("10"))

	idleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("14"))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("11"))
)

type ChatEntry struct {
	From    string
	To      string
	Content string
	IsUser  bool
}

type impulseMsg struct {
	origin  string
	content string
}

type universeMsg struct {
	state reality.UniverseState
}

type Model struct {
	tui           *reality.TUITerm
	width         int
	height        int
	input         textarea.Model
	inputHeight   int
	messages      []ChatEntry
	universe      *reality.UniverseState
	chat          viewport.Model
	ready         bool
	sidebarFocus  bool
	sidebarIdx    int
	selectedBeing string
	lastSpeaker   string
}

const maxInputHeight = 10

func New(t *reality.TUITerm) Model {
	ti := textarea.New()
	ti.Placeholder = "speak..."
	ti.Focus()
	ti.SetHeight(1)
	ti.ShowLineNumbers = false
	ti.KeyMap.InsertNewline.SetKeys("alt+enter")
	ti.CharLimit = 0

	return Model{
		tui:         t,
		input:       ti,
		inputHeight: 1,
		messages:    []ChatEntry{},
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		textarea.Blink,
		listenDisplay(m.tui),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit

		case tea.KeyTab:
			m.sidebarFocus = !m.sidebarFocus
			if m.sidebarFocus {
				m.input.Blur()
			} else {
				m.input.Focus()
			}
			return m, nil

		case tea.KeyEsc:
			if m.sidebarFocus {
				m.sidebarFocus = false
				m.input.Focus()
			} else if m.selectedBeing != "" {
				m.selectedBeing = ""
				m.updateChat()
			}
			return m, nil

		case tea.KeyUp:
			if m.sidebarFocus {
				if m.sidebarIdx > 0 {
					m.sidebarIdx--
				}
				return m, nil
			}

		case tea.KeyDown:
			if m.sidebarFocus {
				max := m.beingCount() - 1
				if max < 0 {
					max = 0
				}
				if m.sidebarIdx < max {
					m.sidebarIdx++
				}
				return m, nil
			}

		case tea.KeyEnter:
			if m.sidebarFocus {
				if name := m.beingAtIdx(m.sidebarIdx); name != "" {
					m.selectedBeing = name
					m.sidebarFocus = false
					m.input.Focus()
					m.updateChat()
				}
				return m, nil
			}

			val := strings.TrimSpace(m.input.Value())
			if val == "" {
				return m, nil
			}

			send := val
			target := m.selectedBeing
			if target != "" && target != m.lastSpeaker {
				send = target + " " + val
			}

			m.messages = append(m.messages, ChatEntry{
				From: "michael", To: target, Content: val, IsUser: true,
			})
			m.input.Reset()
			m.inputHeight = 1
			m.input.SetHeight(1)
			m.resizeChat()
			m.updateChat()
			t := m.tui
			return m, func() tea.Msg {
				t.Send(send)
				return nil
			}
		}

	case tea.MouseMsg:
		if m.ready {
			var cmd tea.Cmd
			m.chat, cmd = m.chat.Update(msg)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		if !m.ready {
			chatW := m.width - sidebarWidth - 6
			chatH := m.chatHeight()
			m.chat = viewport.New(chatW, chatH)
			m.ready = true
		}

		m.input.SetWidth(m.width - 4)
		m.resizeChat()
		m.updateChat()
		return m, nil

	case impulseMsg:
		m.lastSpeaker = msg.origin
		to := "michael"
		m.messages = append(m.messages, ChatEntry{
			From: msg.origin, To: to, Content: msg.content,
		})
		m.updateChat()
		return m, listenDisplay(m.tui)

	case universeMsg:
		m.universe = &msg.state
		m.updateChat()
		return m, listenDisplay(m.tui)
	}

	if !m.sidebarFocus {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		cmds = append(cmds, cmd)

		newHeight := m.calcInputHeight()
		if newHeight != m.inputHeight {
			m.inputHeight = newHeight
			m.input.SetHeight(newHeight)
			m.resizeChat()
			m.updateChat()
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) beingCount() int {
	if m.universe == nil {
		return 0
	}
	count := 0
	for _, b := range m.universe.Beings {
		if b.Name != "michael" {
			count++
		}
	}
	return count
}

func (m *Model) beingAtIdx(idx int) string {
	if m.universe == nil {
		return ""
	}
	i := 0
	for _, b := range m.universe.Beings {
		if b.Name == "michael" {
			continue
		}
		if i == idx {
			return b.Name
		}
		i++
	}
	return ""
}

func (m *Model) calcInputHeight() int {
	val := m.input.Value()
	inputWidth := m.width - 6
	if inputWidth < 1 {
		inputWidth = 1
	}
	visual := 0
	for _, line := range strings.Split(val, "\n") {
		if len(line) == 0 {
			visual++
		} else {
			visual += (len(line)-1)/inputWidth + 1
		}
	}
	if visual < 1 {
		visual = 1
	}
	if visual > maxInputHeight {
		visual = maxInputHeight
	}
	return visual
}

func (m *Model) chatHeight() int {
	h := m.height - m.inputHeight - 6
	if h < 1 {
		h = 1
	}
	return h
}

func (m *Model) resizeChat() {
	if !m.ready {
		return
	}
	m.chat.Width = m.width - sidebarWidth - 6
	m.chat.Height = m.chatHeight()
}

func (m *Model) updateChat() {
	if !m.ready {
		return
	}

	var sb strings.Builder

	if m.selectedBeing != "" && m.universe != nil {
		sb.WriteString(dimStyle.Render(fmt.Sprintf("  exchange: michael ↔ %s", m.selectedBeing)) + "\n\n")
		entries := m.exchangeEntries(m.selectedBeing)
		for _, entry := range entries {
			var name string
			if entry.From == "michael" {
				name = senderUser.Render(entry.From)
			} else {
				name = senderBeing.Render(entry.From)
			}
			lines := wrapText(entry.Content, m.chat.Width-4)
			sb.WriteString(name + "\n")
			for _, line := range lines {
				sb.WriteString("  " + line + "\n")
			}
			sb.WriteString("\n")
		}

		for _, msg := range m.messages {
			if msg.To != m.selectedBeing && msg.From != m.selectedBeing {
				continue
			}
			if m.isInExchangeSnapshot(msg) {
				continue
			}
			var name string
			if msg.IsUser {
				name = senderUser.Render(msg.From)
			} else {
				name = senderBeing.Render(msg.From)
			}
			lines := wrapText(msg.Content, m.chat.Width-4)
			sb.WriteString(name + "\n")
			for _, line := range lines {
				sb.WriteString("  " + line + "\n")
			}
			sb.WriteString("\n")
		}
	} else {
		for _, entry := range m.messages {
			var name string
			if entry.IsUser {
				name = senderUser.Render(entry.From)
			} else {
				name = senderBeing.Render(entry.From)
			}
			lines := wrapText(entry.Content, m.chat.Width-4)
			sb.WriteString(name + "\n")
			for _, line := range lines {
				sb.WriteString("  " + line + "\n")
			}
			sb.WriteString("\n")
		}
	}

	m.chat.SetContent(sb.String())
	m.chat.GotoBottom()
}

func (m *Model) exchangeEntries(being string) []ChatEntry {
	if m.universe == nil {
		return nil
	}
	key := exchangeKey("michael", being)
	for _, ex := range m.universe.Exchanges {
		if ex.Key == key {
			var entries []ChatEntry
			for _, e := range ex.Entries {
				entries = append(entries, ChatEntry{
					From:    e.From,
					Content: e.Content,
					IsUser:  e.From == "michael",
				})
			}
			return entries
		}
	}
	return nil
}

func (m *Model) isInExchangeSnapshot(msg ChatEntry) bool {
	if m.universe == nil || m.selectedBeing == "" {
		return false
	}
	key := exchangeKey("michael", m.selectedBeing)
	for _, ex := range m.universe.Exchanges {
		if ex.Key != key {
			continue
		}
		for _, e := range ex.Entries {
			if e.From == msg.From && e.Content == msg.Content {
				return true
			}
		}
	}
	return false
}

func exchangeKey(a, b string) string {
	if a < b {
		return a + ":" + b
	}
	return b + ":" + a
}

func (m Model) View() string {
	if !m.ready {
		return "loading..."
	}

	sidebar := m.renderSidebar()
	sidebarRendered := sidebarStyle.
		Width(sidebarWidth).
		Height(m.height - 7).
		Render(sidebar)

	chatRendered := chatStyle.
		Width(m.width - sidebarWidth - 6).
		Height(m.height - 7).
		Render(m.chat.View())

	top := lipgloss.JoinHorizontal(lipgloss.Top, sidebarRendered, chatRendered)

	inputRendered := inputStyle.
		Width(m.width - 4).
		Render(m.input.View())

	return lipgloss.JoinVertical(lipgloss.Left,
		titleStyle.Render(" skyra v.05"),
		top,
		inputRendered,
	)
}

func (m Model) renderSidebar() string {
	var sb strings.Builder

	if m.sidebarFocus {
		sb.WriteString(titleStyle.Render("BEINGS") + "\n")
	} else {
		sb.WriteString(dimStyle.Render("BEINGS") + "\n")
	}
	sb.WriteString(dimStyle.Render(strings.Repeat("─", sidebarWidth-4)) + "\n")

	if m.universe != nil {
		idx := 0
		for _, b := range m.universe.Beings {
			if b.Name == "michael" {
				continue
			}

			var status string
			if b.Status == "active" {
				status = activeStyle.Render("●")
			} else {
				status = idleStyle.Render("○")
			}

			name := b.Name
			isCursor := m.sidebarFocus && idx == m.sidebarIdx
			isSelected := m.selectedBeing == b.Name

			if isCursor {
				sb.WriteString(cursorStyle.Render(fmt.Sprintf(" > %s", name)))
			} else if isSelected {
				sb.WriteString(fmt.Sprintf(" %s %s", status, selectedStyle.Render(name)))
			} else {
				sb.WriteString(fmt.Sprintf(" %s %s", status, name))
			}
			sb.WriteString(dimStyle.Render(fmt.Sprintf(" [%s]", b.Type)))
			sb.WriteString("\n")

			if b.Level != nil {
				sb.WriteString(dimStyle.Render(fmt.Sprintf("   Lv.%d  %dxp", b.Level.Level, b.Level.XP)))
				sb.WriteString("\n")
			}
			idx++
		}
	} else {
		sb.WriteString(dimStyle.Render(" waiting...") + "\n")
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("THREADS") + "\n")
	sb.WriteString(dimStyle.Render(strings.Repeat("─", sidebarWidth-4)) + "\n")

	if m.universe != nil {
		active := 0
		for _, t := range m.universe.Threads {
			if t.Active {
				active++
			}
		}
		sb.WriteString(dimStyle.Render(fmt.Sprintf(" Active: %d  Total: %d", active, len(m.universe.Threads))))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("EXCHANGES") + "\n")
	sb.WriteString(dimStyle.Render(strings.Repeat("─", sidebarWidth-4)) + "\n")

	if m.universe != nil {
		active := 0
		for _, e := range m.universe.Exchanges {
			if e.Active {
				active++
			}
		}
		sb.WriteString(dimStyle.Render(fmt.Sprintf(" Active: %d  Total: %d", active, len(m.universe.Exchanges))))
		sb.WriteString("\n")
	}

	sb.WriteString("\n")
	sb.WriteString(dimStyle.Render("MEMORY") + "\n")
	sb.WriteString(dimStyle.Render(strings.Repeat("─", sidebarWidth-4)) + "\n")

	if m.universe != nil {
		for _, b := range m.universe.Beings {
			if b.Type != "llm" {
				continue
			}
			items := len(b.Memories.Items)
			skills := len(b.Memories.Skills)
			if items > 0 || skills > 0 {
				sb.WriteString(fmt.Sprintf(" %s: ", b.Name))
				sb.WriteString(dimStyle.Render(fmt.Sprintf("%d items, %d skills", items, skills)))
				sb.WriteString("\n")
			}
		}
	}

	return sb.String()
}

func listenDisplay(t *reality.TUITerm) tea.Cmd {
	return func() tea.Msg {
		msg := <-t.Display
		switch msg.Type {
		case "impulse":
			return impulseMsg{origin: msg.Origin, content: msg.Content}
		case "universe":
			var state reality.UniverseState
			json.Unmarshal([]byte(msg.Content), &state)
			return universeMsg{state: state}
		}
		return nil
	}
}

func wrapText(text string, width int) []string {
	if width <= 0 {
		width = 80
	}
	var lines []string
	for _, raw := range strings.Split(text, "\n") {
		if len(raw) <= width {
			lines = append(lines, raw)
			continue
		}
		words := strings.Fields(raw)
		var current string
		for _, word := range words {
			if current == "" {
				current = word
			} else if len(current)+1+len(word) <= width {
				current += " " + word
			} else {
				lines = append(lines, current)
				current = word
			}
		}
		if current != "" {
			lines = append(lines, current)
		}
	}
	if len(lines) == 0 {
		lines = []string{""}
	}
	return lines
}

