package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gabrielalmir/vibe-hush/api"
	"github.com/joho/godotenv"
)

func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Parse configuration from environment
	capacity, _ := strconv.Atoi(getEnvWithDefault("VIBE_CACHE_CAPACITY", "1000"))
	expiration, _ := time.ParseDuration(getEnvWithDefault("VIBE_CACHE_EXPIRATION", "10s"))
	port := getEnvWithDefault("VIBE_PORT", "8080")

	// Configure server
	config := api.ServerConfig{
		Capacity:          capacity,
		DefaultExpiration: expiration,
		AuthToken:         os.Getenv("VIBE_AUTH_TOKEN"),
		CertFile:          os.Getenv("VIBE_CERT_FILE"),
		KeyFile:           os.Getenv("VIBE_KEY_FILE"),
	}

	// Create and start server
	server := api.NewCacheServer(config)
	log.Fatal(server.Run(fmt.Sprintf(":%s", port)))
}
