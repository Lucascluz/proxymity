package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"proxymity/internal/config"
	"proxymity/internal/server"
)

func main() {
	// Load config
	cfg, err := config.Load("./config.yaml")
	if err != nil {
		log.Fatal(err)
		return
	}

	// Create server
	srv := server.New(cfg)

	// Start server
	go srv.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	// Give servers 30 seconds to finish on going requests
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited gracefully")
}
