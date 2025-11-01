package controllers

import (
	"github.com/gin-gonic/gin"
)

type Router struct {
	eventHandler *EventHandler
}

func NewRouter(e *EventHandler) *Router {
	return &Router{eventHandler: e}
}

func(r *Router) InitServer() *gin.Engine{
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		events := api.Group("/events")
		{
			events.POST("", r.eventHandler.HandleCreateEvent)
			events.GET("/:id", r.eventHandler.HandleGetEventByID)
			events.GET("", r.eventHandler.HandleListEvents)
		}
	}
	return router
}
