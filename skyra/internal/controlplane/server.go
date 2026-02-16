package controlplane

import (
	"bytes"
	"context"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"
)

type Server struct {
	Orch Orchestrator

	mu       sync.Mutex
	limiters map[string]*ipWindow
}

type ipWindow struct {
	windowStart time.Time
	count       int
}

const (
	maxJSONBodyBytes = 1 << 20 // 1 MiB
	rateLimitWindow  = time.Minute
	rateLimitMaxReq  = 30
	chatTimeout      = 30 * time.Second
	modelRetryCount  = 2
	retryBackoffBase = 200 * time.Millisecond
	defaultTemp      = 0.2
	defaultMaxTokens = 512
	maxMaxTokens     = 4096
)

var thinkBlockRE = regexp.MustCompile(`(?s)<think>.*?</think>`)

// Entry point used by main()
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()

	// core system endpoints
	s.registerSystemRoutes(mux)

	// versioned API
	s.registerV1Routes(mux)

	// wrap everything in middleware
	return s.withMiddleware(mux)
}

//
// Route registration
//

func (s *Server) registerSystemRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/health", s.handleHealth)
}

func (s *Server) registerV1Routes(mux *http.ServeMux) {
	mux.HandleFunc("/v1/chat", s.handleChat)
	mux.HandleFunc("/v1/voice", s.handleVoice)
}

//
// Middleware
//

func (s *Server) withMiddleware(next http.Handler) http.Handler {
	if s.limiters == nil {
		s.limiters = make(map[string]*ipWindow)
	}

	return s.loggingMiddleware(
		s.recoverMiddleware(
			s.rateLimitMiddleware(
				s.authMiddleware(next),
			),
		),
	)
}

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(sw, r)

		log.Printf(
			"%s %s %d %s",
			r.Method,
			r.URL.Path,
			sw.status,
			time.Since(start),
		)
	})
}

func (s *Server) recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/") {
			next.ServeHTTP(w, r)
			return
		}

		expected := strings.TrimSpace(os.Getenv("SKYRA_API_KEY"))
		if expected == "" {
			http.Error(w, "server not configured", http.StatusServiceUnavailable)
			return
		}

		provided := extractAPIKey(r)
		if subtle.ConstantTimeCompare([]byte(expected), []byte(provided)) != 1 {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, "/v1/") {
			next.ServeHTTP(w, r)
			return
		}

		ip := clientIP(r)
		now := time.Now()

		s.mu.Lock()
		win, ok := s.limiters[ip]
		if !ok || now.Sub(win.windowStart) >= rateLimitWindow {
			win = &ipWindow{windowStart: now, count: 0}
			s.limiters[ip] = win
		}
		win.count++
		limited := win.count > rateLimitMaxReq
		s.mu.Unlock()

		if limited {
			w.Header().Set("Retry-After", "60")
			http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

//
// Handlers
//

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (s *Server) handleChat(w http.ResponseWriter, r *http.Request) {
	if !requireJSONPost(w, r) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)

	req, err := decodeChatRequest(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), chatTimeout)
	defer cancel()

	resp, statusCode, err := s.generateChatResponse(ctx, req)
	if err != nil {
		log.Printf("chat generation failed: %v", err)
		http.Error(w, http.StatusText(statusCode), statusCode)
		return
	}

	writeJSON(w, http.StatusOK, resp)
}

func (s *Server) handleVoice(w http.ResponseWriter, r *http.Request) {
	if !requireJSONPost(w, r) {
		return
	}
	r.Body = http.MaxBytesReader(w, r.Body, maxJSONBodyBytes)

	// parse → call orchestrator → write response
	w.Write([]byte("voice endpoint"))
}

func requireJSONPost(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return false
	}

	contentType := strings.ToLower(r.Header.Get("Content-Type"))
	if !strings.HasPrefix(contentType, "application/json") {
		http.Error(w, "content-type must be application/json", http.StatusUnsupportedMediaType)
		return false
	}
	return true
}

func decodeChatRequest(body io.Reader) (chatRequest, error) {
	var req chatRequest

	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return chatRequest{}, fmt.Errorf("invalid JSON body")
	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		return chatRequest{}, fmt.Errorf("body must contain a single JSON object")
	}

	if len(req.Messages) == 0 {
		return chatRequest{}, fmt.Errorf("messages is required")
	}

	for i, msg := range req.Messages {
		role := strings.TrimSpace(msg.Role)
		content := strings.TrimSpace(msg.Content)
		if role == "" || content == "" {
			return chatRequest{}, fmt.Errorf("messages[%d] requires role and content", i)
		}
	}

	if req.MaxTokens != nil {
		if *req.MaxTokens <= 0 {
			return chatRequest{}, fmt.Errorf("max_tokens must be > 0")
		}
		if *req.MaxTokens > maxMaxTokens {
			req.MaxTokens = intPtr(maxMaxTokens)
		}
	}

	if req.Temperature != nil {
		if *req.Temperature < 0 || *req.Temperature > 2 {
			return chatRequest{}, fmt.Errorf("temperature must be between 0 and 2")
		}
	}

	return req, nil
}

func (s *Server) generateChatResponse(ctx context.Context, req chatRequest) (chatResponse, int, error) {
	endpoint := strings.TrimSpace(os.Getenv("SKYRA_MODEL_ENDPOINT"))
	if endpoint == "" {
		return chatResponse{}, http.StatusServiceUnavailable, errors.New("missing SKYRA_MODEL_ENDPOINT")
	}

	model := strings.TrimSpace(req.Model)
	if model == "" {
		model = strings.TrimSpace(os.Getenv("SKYRA_MODEL_NAME"))
	}
	if model == "" {
		return chatResponse{}, http.StatusServiceUnavailable, errors.New("missing model name")
	}

	temperature := defaultTemp
	if req.Temperature != nil {
		temperature = *req.Temperature
	}

	maxTokens := defaultMaxTokens
	if req.MaxTokens != nil {
		maxTokens = *req.MaxTokens
	}

	upReq := upstreamChatRequest{
		Model:       model,
		Messages:    req.Messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
	}

	var lastErr error
	for attempt := 0; attempt <= modelRetryCount; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(1<<uint(attempt-1)) * retryBackoffBase
			select {
			case <-ctx.Done():
				return chatResponse{}, http.StatusGatewayTimeout, ctx.Err()
			case <-time.After(backoff):
			}
		}

		resp, statusCode, err := doUpstreamChatCompletion(ctx, endpoint, upReq)
		if err == nil {
			return resp, http.StatusOK, nil
		}

		lastErr = err
		if statusCode < 500 {
			return chatResponse{}, statusCode, err
		}
	}

	return chatResponse{}, http.StatusBadGateway, lastErr
}

func doUpstreamChatCompletion(ctx context.Context, endpoint string, req upstreamChatRequest) (chatResponse, int, error) {
	payload, err := json.Marshal(req)
	if err != nil {
		return chatResponse{}, http.StatusInternalServerError, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return chatResponse{}, http.StatusInternalServerError, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: chatTimeout}
	httpResp, err := client.Do(httpReq)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			return chatResponse{}, http.StatusGatewayTimeout, err
		}
		return chatResponse{}, http.StatusBadGateway, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(httpResp.Body, 512))
		return chatResponse{}, mapUpstreamErrorStatus(httpResp.StatusCode), fmt.Errorf("model endpoint returned %d: %s", httpResp.StatusCode, strings.TrimSpace(string(body)))
	}

	var upResp upstreamChatResponse
	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(&upResp); err != nil {
		return chatResponse{}, http.StatusBadGateway, fmt.Errorf("invalid model response")
	}

	if len(upResp.Choices) == 0 {
		return chatResponse{}, http.StatusBadGateway, fmt.Errorf("model response missing choices")
	}

	raw := upResp.Choices[0].Message.Content
	sanitized := sanitizeAssistantText(raw)

	return chatResponse{
		Reply: sanitized,
		Model: upResp.Model,
		Usage: upResp.Usage,
	}, http.StatusOK, nil
}

func mapUpstreamErrorStatus(status int) int {
	switch {
	case status == http.StatusUnauthorized || status == http.StatusForbidden:
		return http.StatusBadGateway
	case status == http.StatusTooManyRequests:
		return http.StatusBadGateway
	case status >= 500:
		return http.StatusBadGateway
	default:
		return http.StatusBadGateway
	}
}

func sanitizeAssistantText(raw string) string {
	s := thinkBlockRE.ReplaceAllString(raw, "")
	s = strings.ReplaceAll(s, "</think>", "")
	s = strings.TrimSpace(s)
	s = dedupeRepeatedText(s)
	return strings.TrimSpace(s)
}

func dedupeRepeatedText(s string) string {
	if s == "" {
		return s
	}

	trimmed := strings.TrimSpace(s)

	if idx := strings.Index(trimmed, "\n\n"); idx > 0 {
		left := strings.TrimSpace(trimmed[:idx])
		right := strings.TrimSpace(trimmed[idx+2:])
		if left != "" && left == right {
			return left
		}
	}

	if len(trimmed)%2 == 0 {
		half := len(trimmed) / 2
		left := strings.TrimSpace(trimmed[:half])
		right := strings.TrimSpace(trimmed[half:])
		if left != "" && left == right {
			return left
		}
	}

	return trimmed
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(payload); err != nil {
		log.Printf("write response failed: %v", err)
	}
}

func intPtr(v int) *int {
	return &v
}

func extractAPIKey(r *http.Request) string {
	auth := strings.TrimSpace(r.Header.Get("Authorization"))
	if len(auth) > len("Bearer ") && strings.EqualFold(auth[:len("Bearer ")], "Bearer ") {
		return strings.TrimSpace(auth[len("Bearer "):])
	}
	return strings.TrimSpace(r.Header.Get("X-API-Key"))
}

func clientIP(r *http.Request) string {
	if ip := strings.TrimSpace(r.Header.Get("CF-Connecting-IP")); ip != "" {
		return ip
	}

	if xff := strings.TrimSpace(r.Header.Get("X-Forwarded-For")); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}

	host, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr))
	if err == nil {
		return host
	}
	return strings.TrimSpace(r.RemoteAddr)
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

type chatRequest struct {
	Model       string        `json:"model,omitempty"`
	Messages    []chatMessage `json:"messages"`
	Temperature *float64      `json:"temperature,omitempty"`
	MaxTokens   *int          `json:"max_tokens,omitempty"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatResponse struct {
	Reply string        `json:"reply"`
	Model string        `json:"model,omitempty"`
	Usage upstreamUsage `json:"usage,omitempty"`
}

type upstreamChatRequest struct {
	Model       string        `json:"model"`
	Messages    []chatMessage `json:"messages"`
	MaxTokens   int           `json:"max_tokens"`
	Temperature float64       `json:"temperature"`
}

type upstreamChatResponse struct {
	Model   string           `json:"model"`
	Choices []upstreamChoice `json:"choices"`
	Usage   upstreamUsage    `json:"usage"`
}

type upstreamChoice struct {
	Message chatMessage `json:"message"`
}

type upstreamUsage struct {
	PromptTokens     int `json:"prompt_tokens,omitempty"`
	CompletionTokens int `json:"completion_tokens,omitempty"`
	TotalTokens      int `json:"total_tokens,omitempty"`
}
