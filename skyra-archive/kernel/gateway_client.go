package kernel

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type GatewayClient struct {
	wsURL string
	token string
	mu    sync.Mutex
}

type gatewayEvent struct {
	Type      string          `json:"type"`
	Code      string          `json:"code,omitempty"`
	Delta     string          `json:"delta,omitempty"`
	RequestID string          `json:"request_id,omitempty"`
	SessionID string          `json:"session_id,omitempty"`
	Message   json.RawMessage `json:"message,omitempty"`
}

type gatewayChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func NewGatewayClient(wsURL, token string) *GatewayClient {
	return &GatewayClient{
		wsURL: strings.TrimSpace(wsURL),
		token: strings.TrimSpace(token),
	}
}

func (c *GatewayClient) Ready(ctx context.Context) error {
	if c == nil || c.wsURL == "" {
		return fmt.Errorf("gateway client is not configured")
	}
	return checkHTTPReady(ctx, c.wsURL)
}

func (c *GatewayClient) RunPrompt(ctx context.Context, prompt string, onDelta func(string)) (response string, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	startedAt := time.Now()
	sessionID := "sess_" + newID()
	requestID := "req_" + newID()
	lastPhase := "connect"
	var chatSentAt time.Time
	var firstDeltaAt time.Time
	var streamedChars int
	var sawChatStarted bool

	log.Printf(
		"gateway: request started request=%s session=%s prompt_chars=%d ws_url=%s",
		requestID,
		sessionID,
		len(prompt),
		c.wsURL,
	)
	defer func() {
		status := "ok"
		switch {
		case err == nil:
		case errors.Is(err, context.DeadlineExceeded):
			status = "timeout"
		case errors.Is(err, context.Canceled):
			status = "canceled"
		default:
			status = "error"
		}

		firstToken := "-"
		if !firstDeltaAt.IsZero() && !chatSentAt.IsZero() {
			firstToken = firstDeltaAt.Sub(chatSentAt).Round(time.Millisecond).String()
		}

		modelTime := "-"
		if !chatSentAt.IsZero() {
			modelTime = time.Since(chatSentAt).Round(time.Millisecond).String()
		}

		msg := fmt.Sprintf(
			"gateway: request finished request=%s session=%s status=%s last_phase=%s total=%s model_time=%s first_token=%s streamed_chars=%d response_chars=%d",
			requestID,
			sessionID,
			status,
			lastPhase,
			time.Since(startedAt).Round(time.Millisecond),
			modelTime,
			firstToken,
			streamedChars,
			len(response),
		)
		if err != nil {
			log.Printf("%s err=%v", msg, err)
			return
		}
		log.Print(msg)
	}()

	connectStartedAt := time.Now()
	conn, err := dialWebSocket(ctx, c.wsURL)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	log.Printf(
		"gateway: phase=connect request=%s session=%s duration=%s",
		requestID,
		sessionID,
		time.Since(connectStartedAt).Round(time.Millisecond),
	)

	done := make(chan struct{})
	defer close(done)
	go func() {
		select {
		case <-ctx.Done():
			_ = conn.Close()
		case <-done:
		}
	}()

	lastPhase = "auth.write"
	authStartedAt := time.Now()
	if err := writeGatewayJSON(conn, map[string]any{
		"type":  "auth",
		"token": c.token,
	}); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", err
	}

	lastPhase = "auth.wait"
	if err := waitForEvent(ctx, conn, func(evt gatewayEvent) (bool, error) {
		switch evt.Type {
		case "auth.ok":
			return true, nil
		case "error":
			return false, fmt.Errorf("gateway auth failed: %s", evt.Code)
		default:
			return false, nil
		}
	}); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", err
	}
	log.Printf(
		"gateway: phase=auth request=%s session=%s duration=%s",
		requestID,
		sessionID,
		time.Since(authStartedAt).Round(time.Millisecond),
	)

	lastPhase = "session.open.write"
	sessionOpenStartedAt := time.Now()
	if err := writeGatewayJSON(conn, map[string]any{
		"type":       "session.open",
		"session_id": sessionID,
	}); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", err
	}

	lastPhase = "session.open.wait"
	if err := waitForEvent(ctx, conn, func(evt gatewayEvent) (bool, error) {
		switch evt.Type {
		case "session.opened":
			if evt.SessionID == sessionID {
				return true, nil
			}
			return false, nil
		case "warning":
			return false, nil
		case "error":
			return false, fmt.Errorf("gateway session error: %s", evt.Code)
		default:
			return false, nil
		}
	}); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", err
	}
	log.Printf(
		"gateway: phase=session.open request=%s session=%s duration=%s",
		requestID,
		sessionID,
		time.Since(sessionOpenStartedAt).Round(time.Millisecond),
	)

	lastPhase = "chat.send"
	chatSentAt = time.Now()
	if err := writeGatewayJSON(conn, map[string]any{
		"type":       "chat.send",
		"request_id": requestID,
		"session_id": sessionID,
		"messages": []gatewayChatMessage{
			{Role: "user", Content: prompt},
		},
		"stream": true,
	}); err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		return "", err
	}
	log.Printf(
		"gateway: phase=chat.send request=%s session=%s prompt_chars=%d",
		requestID,
		sessionID,
		len(prompt),
	)

	pingCtx, cancelPing := context.WithCancel(context.Background())
	defer cancelPing()
	go gatewayPingLoop(pingCtx, conn)

	var builder strings.Builder
	for {
		raw, err := conn.ReadText(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return "", ctx.Err()
			}
			return "", err
		}

		var evt gatewayEvent
		if err := json.Unmarshal([]byte(raw), &evt); err != nil {
			return "", fmt.Errorf("decode gateway event: %w", err)
		}

		switch evt.Type {
		case "pong":
			continue
		case "chat.started":
			if evt.RequestID != "" && evt.RequestID != requestID {
				continue
			}
			if !sawChatStarted {
				sawChatStarted = true
				lastPhase = "chat.started"
				log.Printf(
					"gateway: phase=chat.started request=%s session=%s queued=%s",
					requestID,
					sessionID,
					time.Since(chatSentAt).Round(time.Millisecond),
				)
			}
			continue
		case "warning":
			log.Printf(
				"gateway: phase=warning request=%s session=%s code=%s",
				requestID,
				sessionID,
				evt.Code,
			)
			continue
		case "chat.delta":
			if evt.RequestID != requestID {
				continue
			}
			if firstDeltaAt.IsZero() {
				firstDeltaAt = time.Now()
				lastPhase = "chat.delta"
				log.Printf(
					"gateway: phase=first_token request=%s session=%s after_send=%s total=%s",
					requestID,
					sessionID,
					firstDeltaAt.Sub(chatSentAt).Round(time.Millisecond),
					firstDeltaAt.Sub(startedAt).Round(time.Millisecond),
				)
			}
			builder.WriteString(evt.Delta)
			streamedChars += len(evt.Delta)
			if onDelta != nil && evt.Delta != "" {
				onDelta(evt.Delta)
			}
		case "chat.done":
			if evt.RequestID != requestID {
				continue
			}
			lastPhase = "chat.done"
			var message gatewayChatMessage
			if len(evt.Message) != 0 {
				if err := json.Unmarshal(evt.Message, &message); err != nil {
					return "", err
				}
			}
			if message.Content != "" {
				response = message.Content
			} else {
				response = builder.String()
			}
			log.Printf(
				"gateway: phase=chat.done request=%s session=%s model_time=%s streamed_chars=%d response_chars=%d",
				requestID,
				sessionID,
				time.Since(chatSentAt).Round(time.Millisecond),
				streamedChars,
				len(response),
			)
			return response, nil
		case "chat.canceled":
			if evt.RequestID == requestID {
				return "", fmt.Errorf("gateway request canceled")
			}
		case "error":
			if evt.RequestID == "" || evt.RequestID == requestID {
				return "", fmt.Errorf("gateway error %s", evt.Code)
			}
		}
	}
}

func gatewayPingLoop(ctx context.Context, conn *wsClient) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case t := <-ticker.C:
			_ = writeGatewayJSON(conn, map[string]any{
				"type": "ping",
				"ts":   t.UnixMilli(),
			})
		}
	}
}

func waitForEvent(ctx context.Context, conn *wsClient, accept func(gatewayEvent) (bool, error)) error {
	for {
		raw, err := conn.ReadText(ctx)
		if err != nil {
			return err
		}
		var evt gatewayEvent
		if err := json.Unmarshal([]byte(raw), &evt); err != nil {
			return err
		}
		done, err := accept(evt)
		if err != nil {
			return err
		}
		if done {
			return nil
		}
	}
}

func writeGatewayJSON(conn *wsClient, payload any) error {
	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	return conn.WriteText(string(buf))
}
