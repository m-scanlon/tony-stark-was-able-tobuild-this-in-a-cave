package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/mikescanlon/skyra/kernel"
)

const (
	defaultAddr          = "127.0.0.1:9090"
	defaultValkeyAddr    = "127.0.0.1:6379"
	defaultHeapCapacity  = 1024
	maxIngressBodyBytes  = 1 << 20 // 1 MiB
	bootstrapSkillPrefix = "skill:"
)

type config struct {
	Addr           string
	ValkeyAddr     string
	ValkeyUsername string
	ValkeyPassword string
	HeapCapacity   int
	BootstrapTools []string
}

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("startup: %v", err)
	}

	valkeyClient := kernel.NewValkeyClient(cfg.ValkeyAddr, cfg.ValkeyUsername, cfg.ValkeyPassword)

	startupCtx, startupCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer startupCancel()
	if err := valkeyClient.Ping(startupCtx); err != nil {
		log.Fatalf("startup: valkey unavailable at %q: %v", cfg.ValkeyAddr, err)
	}

	if err := seedBootstrapSkills(startupCtx, valkeyClient, cfg.BootstrapTools); err != nil {
		log.Fatalf("startup: failed to seed bootstrap skills: %v", err)
	}

	registry := kernel.NewValkeyRegistry(valkeyClient)
	heap := kernel.NewInMemoryHeap(cfg.HeapCapacity)
	k := kernel.New(registry, heap)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go k.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/v1/internal/health", handleHealth)
	mux.HandleFunc("/dispatch", handleIngress(k, ctx)) // compatibility path
	mux.HandleFunc("/v1/internal/ingress", handleIngress(k, ctx))
	mux.HandleFunc("/v1/internal/egress", handleEgress())

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("internal api listening on %s", cfg.Addr)
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
		cancel()
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			log.Printf("shutdown failed: %v", err)
			os.Exit(1)
		}
	case err := <-errCh:
		if err != nil {
			log.Printf("server failed: %v", err)
			os.Exit(1)
		}
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type ingressRequest struct {
	Command string `json:"command"`
	TurnID  string `json:"turn_id,omitempty"`
}

type ingressResponse struct {
	OK         bool   `json:"ok"`
	AckCommand string `json:"ack_command,omitempty"`
	Error      string `json:"error,omitempty"`
}

func handleIngress(k *kernel.Kernel, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req, err := parseIngressRequest(w, r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ingressResponse{Error: err.Error()})
			return
		}

		if err := k.Dispatch(ctx, req.Command); err != nil {
			writeJSON(w, http.StatusUnprocessableEntity, ingressResponse{Error: err.Error()})
			return
		}

		resp := ingressResponse{OK: true}
		if req.TurnID != "" {
			resp.AckCommand = fmt.Sprintf("octos ack --turn=%s --status=stored", req.TurnID)
		}
		writeJSON(w, http.StatusOK, resp)
	}
}

func handleEgress() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		req, err := parseIngressRequest(w, r)
		if err != nil {
			writeJSON(w, http.StatusBadRequest, ingressResponse{Error: err.Error()})
			return
		}

		if _, err := kernel.ParseCommand(req.Command); err != nil {
			writeJSON(w, http.StatusUnprocessableEntity, ingressResponse{Error: err.Error()})
			return
		}

		log.Printf("egress command accepted: %s", req.Command)
		writeJSON(w, http.StatusOK, ingressResponse{OK: true})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("writeJSON: encode failed: %v", err)
	}
}

func parseIngressRequest(w http.ResponseWriter, r *http.Request) (ingressRequest, error) {
	r.Body = http.MaxBytesReader(w, r.Body, maxIngressBodyBytes)
	defer r.Body.Close()

	contentType := strings.ToLower(strings.TrimSpace(r.Header.Get("Content-Type")))
	if strings.HasPrefix(contentType, "text/plain") {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return ingressRequest{}, errors.New("invalid request body")
		}
		command := strings.TrimSpace(string(body))
		if command == "" {
			return ingressRequest{}, errors.New("command is required")
		}
		return ingressRequest{Command: command}, nil
	}

	var req ingressRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		return ingressRequest{}, errors.New("invalid request body")
	}
	req.Command = strings.TrimSpace(req.Command)
	req.TurnID = strings.TrimSpace(req.TurnID)
	if req.Command == "" {
		return ingressRequest{}, errors.New("command is required")
	}

	return req, nil
}

func loadConfig() (config, error) {
	addr := strings.TrimSpace(firstNonEmpty(os.Getenv("INTERNAL_API_ADDR"), os.Getenv("KERNEL_ADDR"), defaultAddr))
	valkeyAddr := strings.TrimSpace(firstNonEmpty(os.Getenv("VALKEY_ADDR"), defaultValkeyAddr))
	valkeyUsername := strings.TrimSpace(os.Getenv("VALKEY_USERNAME"))
	valkeyPassword := strings.TrimSpace(os.Getenv("VALKEY_PASSWORD"))
	if valkeyPassword == "" && os.Getenv("VALKEY_ALLOW_NOAUTH") != "1" {
		return config{}, errors.New("VALKEY_PASSWORD is required (set VALKEY_ALLOW_NOAUTH=1 for local unsecured mode)")
	}

	heapCapacity := defaultHeapCapacity
	if raw := strings.TrimSpace(os.Getenv("KERNEL_HEAP_CAPACITY")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			return config{}, fmt.Errorf("invalid KERNEL_HEAP_CAPACITY: %q", raw)
		}
		heapCapacity = n
	}

	return config{
		Addr:           addr,
		ValkeyAddr:     valkeyAddr,
		ValkeyUsername: valkeyUsername,
		ValkeyPassword: valkeyPassword,
		HeapCapacity:   heapCapacity,
		BootstrapTools: parseCSVTools(os.Getenv("INTERNAL_API_SEED_SKILLS")),
	}, nil
}

func seedBootstrapSkills(ctx context.Context, client *kernel.ValkeyClient, tools []string) error {
	type seedRecord struct {
		ID          string             `json:"id"`
		Name        string             `json:"name"`
		Description string             `json:"description"`
		Shard       string             `json:"shard"`
		Tasks       []kernel.SkillTask `json:"tasks"`
	}

	for _, tool := range tools {
		key := bootstrapSkillPrefix + tool

		exists, err := client.Exists(ctx, key)
		if err != nil {
			return fmt.Errorf("exists %q: %w", key, err)
		}
		if exists {
			continue
		}

		payload, err := json.Marshal(seedRecord{
			ID:          "skill." + tool,
			Name:        tool,
			Description: "bootstrap skill",
			Shard:       "brain",
			Tasks: []kernel.SkillTask{
				{
					Name:        tool,
					Description: "bootstrap task",
				},
			},
		})
		if err != nil {
			return fmt.Errorf("marshal bootstrap skill %q: %w", tool, err)
		}

		if err := client.Set(ctx, key, string(payload)); err != nil {
			return fmt.Errorf("set %q: %w", key, err)
		}
		log.Printf("seeded bootstrap skill: %s", tool)
	}

	return nil
}

func parseCSVTools(raw string) []string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}

	seen := make(map[string]struct{})
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		tool := strings.TrimSpace(strings.ToLower(part))
		if tool == "" {
			continue
		}
		if _, ok := seen[tool]; ok {
			continue
		}
		seen[tool] = struct{}{}
		out = append(out, tool)
	}

	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return value
		}
	}
	return ""
}
