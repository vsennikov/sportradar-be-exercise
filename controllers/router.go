package controllers

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	eventHandler *EventHandler
	SportHandler *SportHandler
}

func NewRouter(e *EventHandler, s *SportHandler) *Router {
	return &Router{eventHandler: e, SportHandler: s}
}

func(r *Router) InitServer() *gin.Engine{
	router := gin.Default()

	router.Static("/static", "./static")
	router.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})
	api := router.Group("/api/v1")
	{
		sports := api.Group("sports")
		{
			sports.POST("", r.SportHandler.HandleCreateSport)
			sports.GET("/:id", r.SportHandler.HandleGetSportByID)
			sports.GET("", r.SportHandler.HandleListSports)
			sports.DELETE("/:id", r.SportHandler.HandleDeleteSport)
			sports.PUT("/:id", r.SportHandler.HandleUpdateSport)
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
