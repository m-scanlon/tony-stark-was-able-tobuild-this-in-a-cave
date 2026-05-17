package reality

type TUITerm struct {
	id      string
	Device  Reality
	inbox   chan string
	Display chan TUIMessage
}

type TUIMessage struct {
	Type    string // "impulse", "universe"
	Origin  string
	Content string
}

func (t *TUITerm) ID() string { return t.id }

func (t *TUITerm) Create(r *Relation) Reality {
	return &TUITerm{
		id:      "terminal",
		inbox:   make(chan string),
		Display: make(chan TUIMessage, 32),
	}
}

func (t *TUITerm) Realize(r *Relation) string {
	if r.Collecting {
		return ""
	}
	if r.Impulse != "" {
		select {
		case t.Display <- TUIMessage{Type: "impulse", Origin: r.Origin, Content: r.Impulse}:
		default:
		}
	}
	return <-t.inbox
}

func (t *TUITerm) Send(input string) {
	t.inbox <- input
}

func (t *TUITerm) Broadcast(state string) {
	select {
	case t.Display <- TUIMessage{Type: "universe", Content: state}:
	default:
	}
}
