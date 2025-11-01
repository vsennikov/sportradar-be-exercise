package controllers

import "github.com/vsennikov/sportradar-be-exercise/services"

type EventHandler struct {
	eventService services.EventServiceInterface
}

func NewEventHandler(eventService services.EventServiceInterface) *EventHandler {
	return &EventHandler{eventService: eventService}
}