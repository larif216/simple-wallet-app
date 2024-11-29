package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"simple-wallet-app/internal/config"
	"syscall"
	"time"
)

func main() {
	// Initialize HTTP Server
	httpServer, err := config.NewHttpServer()
	if err != nil {
		log.Fatalf("Error initializing the HTTP server: %v", err)
	}

	// Create a channel to listen for termination signals (Ctrl+C, SIGTERM, etc.)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	// Log that the server is starting
	log.Printf("Server starting on %s", httpServer.Config.Host)

	// Run the server in a goroutine so that it doesn't block
	go func() {
		if err := httpServer.HTTPServer.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for a termination signal
	<-stop

	// Graceful shutdown
	log.Println("Shutting down the server...")

	// Create a context for the shutdown to allow pending requests to finish
	shutdownTimeout := 5 * time.Second
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Attempt to shut down the server gracefully
	if err := httpServer.HTTPServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("Error during server shutdown: %v", err)
	}

	log.Println("Server stopped")
}
