package controllers

import (
	"net/http"

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
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return router
}