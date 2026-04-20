package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/yhshin0/go-auth-server/internal/config"
	"github.com/yhshin0/go-auth-server/internal/infrastructure/database"
	"github.com/yhshin0/go-auth-server/internal/middleware"
	"github.com/yhshin0/go-auth-server/internal/router"
)

func main() {
	config.Setup()
	db, err := database.NewDatabase(&config.GetInstance().DB)
	if err != nil {
		log.Fatalf("failed to initialize database: %s\n", err.Error())
	}
	defer db.CloseWithLog()

	// The HTTP Server
	c := config.GetInstance()
	addr := c.Server.Host + ":" + c.Server.Port
	server := &http.Server{
		Addr:              addr,
		Handler:           service(),
		ReadTimeout:       c.Server.HttpReadTimeout,
		ReadHeaderTimeout: c.Server.HttpReadHeaderTimeout,
		WriteTimeout:      c.Server.HttpWriteTimeout,
		IdleTimeout:       c.Server.HttpIdleTimeout,
	}

	// Create context that listens for the interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Run server in the background
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal(err)
		}
	}()

	// Listen for the interrupt signal
	<-ctx.Done()

	// Create shutdown context with 30-second timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Trigger graceful shutdown
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal(err)
	}
}

func service() http.Handler {
	r := chi.NewRouter()

	middleware.Register(r)
	router.Register(r)

	return r
}
