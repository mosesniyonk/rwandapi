package main

import (
	"log"
	"net/http"
	"os"

	"github.com/mosesniyonk/rwandapi/internal/database"
	"github.com/mosesniyonk/rwandapi/internal/server"
)

func main() {
	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Create server
	handler := server.New()

	// Determine port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("rwandapi server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
