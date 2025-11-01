package main

import (
	"log"

	"github.com/vsennikov/sportradar-be-exercise/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}
	router, db, err := buildApp(cfg)
	if err != nil {
		log.Fatalf("Could not build application: %v", err)
	}
	defer cleanupFunc(db)
	log.Printf("Server starting on http://localhost:%s", cfg.AppPort)
	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("FATAL: Could not start server: %v", err)
	}
}
