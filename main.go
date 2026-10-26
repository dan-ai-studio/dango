package main

import (
	"log"
	"os"

	"github.com/dan-ai-studio/dango/internal/server"
)

func main() {
	r := server.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 dango (Gin+Vue) listening on http://localhost:%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
