package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"skyra/internal/controlplane"
)

const (
	defaultAddr          = ":8080"
	readHeaderTimeout    = 5 * time.Second
	readTimeout          = 15 * time.Second
	writeTimeout         = 30 * time.Second
	idleTimeout          = 60 * time.Second
	shutdownGraceTimeout = 10 * time.Second
)

func main() {
	addr := os.Getenv("SKYRA_ADDR")
	if addr == "" {
		addr = defaultAddr
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           (&controlplane.Server{}).Routes(),
		ReadHeaderTimeout: readHeaderTimeout,
		ReadTimeout:       readTimeout,
		WriteTimeout:      writeTimeout,
		IdleTimeout:       idleTimeout,
	}

	errCh := make(chan error, 1)
	go func() {
		log.Printf("skyrad listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
		close(errCh)
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		log.Printf("shutdown signal received: %s", sig)
		ctx, cancel := context.WithTimeout(context.Background(), shutdownGraceTimeout)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			log.Printf("graceful shutdown failed: %v", err)
			os.Exit(1)
		}
		log.Println("server stopped")
	case err := <-errCh:
		if err != nil {
			log.Printf("server failed: %v", err)
			os.Exit(1)
		}
	}
}
