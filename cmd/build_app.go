package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sportradar-be-exercise/config"
	"github.com/vsennikov/sportradar-be-exercise/controllers"
	"github.com/vsennikov/sportradar-be-exercise/infrastructure"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

func cleanupFunc (db *sqlx.DB) {
	if db == nil {
		return
	}
	log.Println("Closing database connection...")
	if err := db.Close(); err != nil {
		log.Printf("Error while closing database: %v", err)
	}
}

func buildApp(cfg config.Config) (*gin.Engine, *sqlx.DB, error) {
	db, err := infrastructure.NewConnection(cfg)
	if err != nil {
		return nil, nil, fmt.Errorf("could not connect to database: %w", err)
	}
	log.Println("Initializing dependencies...")
	eventRepository := infrastructure.NewEventRepository(db)
	sportRepository := infrastructure.NewSportRepository(db)
	venueRepository := infrastructure.NewVenueRepository(db)
	eventService := services.NewEventService(
		eventRepository,
		cfg.DefaultPage,
		cfg.DefaultLimit,
	)
	sportService := services.NewSportService(
		sportRepository,
		eventRepository,
	)
	venueService := services.NewVenueService(
		venueRepository,
		eventRepository,
	)
	sportHandler := controllers.NewSportHandler(sportService)
	eventHandler := controllers.NewEventHandler(eventService)
	venueHandler := controllers.NewVenueHandler(venueService)
	log.Println("Setting up routes...")
	router := controllers.NewRouter(eventHandler, sportHandler, venueHandler)
	server := router.InitServer()
	return server, db, nil
}
