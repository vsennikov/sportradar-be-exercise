package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vsennikov/sportradar-be-exercise/config"
	"github.com/vsennikov/sportradar-be-exercise/infrastructure"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
    log.Fatalf("Failed to parse configuration: %v", err)
	}

	db, err := infrastructure.NewConnection(cfg)
	if err != nil {
		log.Fatalf("%v", err)
	}
	defer db.Close()
	
	r := gin.Default()
	r.GET("/api", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})
	log.Printf("Starting server on port %s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
		if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}