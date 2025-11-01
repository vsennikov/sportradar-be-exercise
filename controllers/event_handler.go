package controllers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type EventHandler struct {
	eventService services.EventServiceInterface
}

func NewEventHandler(eventService services.EventServiceInterface) *EventHandler {
	return &EventHandler{eventService: eventService}
}

func (h *EventHandler) HandleCreateEvent(c *gin.Context) {
	var req services.EventCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newID, err := h.eventService.CreateEvent(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

func (h *EventHandler) HandleGetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event ID format"})
		return
	}
	event, err := h.eventService.GetEventByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	eventDTO := toDTOEvent(*event)
	c.JSON(http.StatusOK, eventDTO)
}

func (h *EventHandler) HandleListEvents(c *gin.Context) {
	var req services.ListEventsRequest

	page, _ := strconv.Atoi(c.Query("page"))
	req.Page = page
	limit, _ := strconv.Atoi(c.Query("limit"))
	req.Limit = limit
	if sportIDStr := c.Query("sport_id"); sportIDStr != "" {
		if sportID, err := strconv.Atoi(sportIDStr); err == nil {
			req.SportID = &sportID
		}
	}
	date_from := c.Query("date_from")
	if (date_from != "") {
		parsedTime, err := time.Parse("2006-01-02", date_from)
		if (err == nil) {
			req.DateFrom = &parsedTime
		}
	}

	events, pagination, err := h.eventService.ListEvents(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	eventDTOs := make([]EventDTO, 0, len(events))
	for _, e := range events {
		eventDTOs = append(eventDTOs, toDTOEvent(e))
	}
	c.JSON(http.StatusOK, gin.H{
		"pagination": pagination,
		"events":     eventDTOs,
	})
}
