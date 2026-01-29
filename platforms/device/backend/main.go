package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"pansiot-device/internal/storage"
)

func main() {
	// Parse command line arguments
	if len(os.Args) < 3 {
		fmt.Println("Usage: pansiot-device -config <config-file>")
		fmt.Println("Example:")
		fmt.Println("  Gateway mode: pansiot-device -config ../../config/gateway.yaml")
		fmt.Println("  HMI mode:     pansiot-device -config ../../config/hmi.yaml")
		os.Exit(1)
	}

	configPath := os.Args[2]
	absConfigPath, err := filepath.Abs(configPath)
	if err != nil {
		log.Fatalf("Failed to resolve config path: %v", err)
	}

	// Load configuration
	log.Printf("Loading configuration from: %s", absConfigPath)
	// TODO: Implement config loading
	// cfg, err := config.Load(absConfigPath)
	// if err != nil {
	// 	log.Fatalf("Failed to load config: %v", err)
	// }

	// Initialize storage layer
	log.Println("Initializing real-time storage layer...")
	storageLayer := storage.NewMemoryStorage()

	// Start the runtime
	log.Println("Starting PansIot Device Platform...")
	// TODO: Start adapters, collectors, consumers based on config

	// Start WebSocket server
	log.Println("Starting WebSocket server...")
	// TODO: Start WebSocket server on configured port

	log.Println("Device platform started successfully")

	// Keep the application running
	select {}
}
