package controllers

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	eventHandler *EventHandler
	sportHandler *SportHandler
	venueHandler *VenueHandler
	teamHandler *TeamHandler
}

func NewRouter(e *EventHandler, s *SportHandler, v *VenueHandler, t *TeamHandler) *Router {
	return &Router{eventHandler: e, sportHandler: s, venueHandler: v, teamHandler: t}
}

func(r *Router) InitServer() *gin.Engine{
	router := gin.Default()

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	api := router.Group("/api/v1")
	{
		teams := api.Group("teams")
		{
			teams.POST("", r.teamHandler.HandleCreateTeam)
			teams.GET("/:id", r.teamHandler.HandleGetTeamByID)
			teams.GET("", r.teamHandler.HandleListTeams)
			teams.PATCH("/:id", r.teamHandler.HandleUpdateTeam)
			teams.DELETE("/:id", r.teamHandler.HandleDeleteTeam)
		}
		venues := api.Group("venues")
		{
			venues.POST("", r.venueHandler.HandleCreateVenue)
			venues.GET("/:id", r.venueHandler.HandleGetVenueByID)
			venues.GET("", r.venueHandler.HandleListVenues)
			venues.PATCH("/:id", r.venueHandler.HandleUpdateVenue)
			venues.DELETE("/:id", r.venueHandler.HandleDeleteVenue)
		}
		sports := api.Group("sports")
		{
			sports.POST("", r.sportHandler.HandleCreateSport)
			sports.GET("/:id", r.sportHandler.HandleGetSportByID)
			sports.GET("", r.sportHandler.HandleListSports)
			sports.DELETE("/:id", r.sportHandler.HandleDeleteSport)
			sports.PUT("/:id", r.sportHandler.HandleUpdateSport)
		}
		events := api.Group("/events")
		{
			events.POST("", r.eventHandler.HandleCreateEvent)
			events.GET("/:id", r.eventHandler.HandleGetEventByID)
			events.GET("", r.eventHandler.HandleListEvents)
		}
	}
	return router
}
