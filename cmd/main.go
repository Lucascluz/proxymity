package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	log.Println(cfg)

	// Create server
	srv := server.New(cfg)

	// Start server
	go srv.Start()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down...")
	srv.Shutdown(context.Background())
}
