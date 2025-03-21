package main

import (
	"log"
)

func main() {
	// Create a new HTTP server on port 1000
	httpServer := NewHttpServer(":8000")

	// Run the server and handle errors
	if err := httpServer.Run(); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}
