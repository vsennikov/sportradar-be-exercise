package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sports-event-calendar/config"
	"github.com/vsennikov/sports-event-calendar/controllers"
	"github.com/vsennikov/sports-event-calendar/infrastructure"
	"github.com/vsennikov/sports-event-calendar/services"
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
	teamRepository := infrastructure.NewTeamRepository(db)
	eventService := services.NewEventService(
		eventRepository,
		cfg.DefaultPage,
		cfg.DefaultLimit,
		sportRepository,
		teamRepository,
		venueRepository,
	)
	sportService := services.NewSportService(
		sportRepository,
		eventRepository,
	)
	venueService := services.NewVenueService(
		venueRepository,
		eventRepository,
	)
	teamService := services.NewTeamService(
		teamRepository,
		eventRepository,
	)
	sportHandler := controllers.NewSportHandler(sportService)
	eventHandler := controllers.NewEventHandler(eventService)
	venueHandler := controllers.NewVenueHandler(venueService)
	teamHandler := controllers.NewTeamHandler(teamService)
	log.Println("Setting up routes...")
	router := controllers.NewRouter(eventHandler, sportHandler, venueHandler, teamHandler)
	server := router.InitServer()
	return server, db, nil
}
