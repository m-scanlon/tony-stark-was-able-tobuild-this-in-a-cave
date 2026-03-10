package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mikescanlon/skyra/kernel"
)

const defaultAddr = ":9090"

func main() {
	addr := os.Getenv("KERNEL_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	// TODO: wire real registry (Redis) and real heap
	k := kernel.New(nil, nil)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go k.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handleHealth)
	mux.HandleFunc("/dispatch", handleDispatch(k, ctx))

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("kernel listening on %s", addr)
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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

type dispatchRequest struct {
	Command string `json:"command"`
}

type dispatchResponse struct {
	OK    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

func handleDispatch(k *kernel.Kernel, ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req dispatchRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, dispatchResponse{Error: "invalid request body"})
			return
		}

		if err := k.Dispatch(ctx, req.Command); err != nil {
			writeJSON(w, http.StatusUnprocessableEntity, dispatchResponse{Error: err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, dispatchResponse{OK: true})
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
