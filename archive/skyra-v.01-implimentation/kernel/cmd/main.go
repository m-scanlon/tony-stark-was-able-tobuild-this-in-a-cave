package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/mikescanlon/skyra/kernel"
)

const (
	defaultAddr         = "0.0.0.0:9090"
	defaultHeapCapacity = 1024
)

//go:embed ui/*
var uiFiles embed.FS

type config struct {
	Addr         string
	HeapCapacity int
	PromptRoot   string
	GatewayWSURL string
	GatewayToken string
}

type stimulusRequest struct {
	Content string `json:"content"`
	Source  string `json:"source,omitempty"`
	Type    string `json:"type,omitempty"`
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	queue := kernel.NewInMemoryHeap(cfg.HeapCapacity)
	k, err := kernel.New(kernel.Config{
		PromptRoot:   cfg.PromptRoot,
		GatewayWSURL: cfg.GatewayWSURL,
		GatewayToken: cfg.GatewayToken,
	}, queue)
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go k.Run(ctx)

	uiSub, err := fs.Sub(uiFiles, "ui")
	if err != nil {
		log.Fatalf("startup: ui filesystem: %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/interact", http.StatusFound))
	mux.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.FS(uiSub))))
	mux.HandleFunc("/interact", serveUI("interact.html"))
	mux.HandleFunc("/thoughts", serveUI("thoughts.html"))
	mux.HandleFunc("/healthz", handleHealth(k))
	mux.HandleFunc("/readyz", handleReady(k))
	mux.HandleFunc("/api/v1/stimuli", handleStimuli(k, ctx))
	mux.HandleFunc("/api/v1/interact", handleInteractSnapshot(k))
	mux.HandleFunc("/api/v1/thoughts", handleThoughtSnapshot(k))
	mux.HandleFunc("/api/v1/interact/stream", handleSSE(k.SubscribeInteraction, k.UnsubscribeInteraction))
	mux.HandleFunc("/api/v1/thoughts/stream", handleSSE(k.SubscribeThoughts, k.UnsubscribeThoughts))

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      0,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("kernel listening on %s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("shutdown signal: %s", sig)
	case err := <-errCh:
		if err != nil {
			log.Printf("server failed: %v", err)
			os.Exit(1)
		}
	}

	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("shutdown failed: %v", err)
		os.Exit(1)
	}
}

func loadConfig() (config, error) {
	addr := firstNonEmpty(os.Getenv("KERNEL_ADDR"), os.Getenv("SKYRA_ADDR"), defaultAddr)
	promptRoot, err := findPromptRoot(os.Getenv("SKYRA_PROMPT_ROOT"))
	if err != nil {
		return config{}, err
	}

	heapCapacity := defaultHeapCapacity
	if raw := strings.TrimSpace(os.Getenv("KERNEL_HEAP_CAPACITY")); raw != "" {
		var parsed int
		if _, err := fmt.Sscanf(raw, "%d", &parsed); err != nil || parsed <= 0 {
			return config{}, fmt.Errorf("invalid KERNEL_HEAP_CAPACITY: %q", raw)
		}
		heapCapacity = parsed
	}

	wsURL := strings.TrimSpace(os.Getenv("OLLAMA_GATEWAY_WS_URL"))
	if wsURL == "" {
		return config{}, errors.New("OLLAMA_GATEWAY_WS_URL is required")
	}

	return config{
		Addr:         addr,
		HeapCapacity: heapCapacity,
		PromptRoot:   promptRoot,
		GatewayWSURL: wsURL,
		GatewayToken: strings.TrimSpace(os.Getenv("OLLAMA_GATEWAY_TOKEN")),
	}, nil
}

func findPromptRoot(explicit string) (string, error) {
	candidates := []string{}
	if strings.TrimSpace(explicit) != "" {
		candidates = append(candidates, explicit)
	}
	candidates = append(candidates,
		"/app/skyra/chain-of-thought/primtives",
		filepath.Join("skyra", "chain-of-thought", "primtives"),
		filepath.Join("..", "chain-of-thought", "primtives"),
	)

	for _, candidate := range candidates {
		if candidate == "" {
			continue
		}
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("could not locate prompt root")
}

func serveUI(name string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		http.ServeFileFS(w, r, uiFiles, filepath.Join("ui", name))
	}
}

func handleHealth(k *kernel.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":          true,
			"service":     "skyra-kernel",
			"queue_depth": k.ThoughtSnapshot().QueueDepth,
		})
	}
}

func handleReady(k *kernel.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()
		if err := k.Ready(ctx); err != nil {
			writeJSON(w, http.StatusServiceUnavailable, map[string]any{
				"ok":      false,
				"message": err.Error(),
			})
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":      true,
			"message": "ready",
		})
	}
}

func handleStimuli(k *kernel.Kernel, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req stimulusRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		content := strings.TrimSpace(req.Content)
		if content == "" {
			http.Error(w, "content is required", http.StatusBadRequest)
			return
		}

		source := strings.TrimSpace(req.Source)
		if source == "" {
			source = "human"
		}

		stimulusType := strings.TrimSpace(req.Type)
		if stimulusType == "" {
			stimulusType = "text"
		}

		stimulus, err := k.SubmitStimulus(ctx, source, stimulusType, content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		writeJSON(w, http.StatusAccepted, map[string]any{
			"ok":       true,
			"stimulus": stimulus,
		})
	}
}

func handleInteractSnapshot(k *kernel.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, k.InteractionSnapshot())
	}
}

func handleThoughtSnapshot(k *kernel.Kernel) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, k.ThoughtSnapshot())
	}
}

func handleSSE(subscribe func() (int, <-chan kernel.UIEvent), unsubscribe func(int)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "streaming unsupported", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		id, ch := subscribe()
		defer unsubscribe(id)

		ctx := r.Context()
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-ch:
				if !ok {
					return
				}
				data, err := json.Marshal(event.Data)
				if err != nil {
					continue
				}
				fmt.Fprintf(w, "event: %s\n", event.Type)
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			}
		}
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
