package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vsennikov/sports-event-calendar/services"
)

type VenueHandler struct {
	venueService services.VenueServiceInterface
}

func NewVenueHandler(s services.VenueServiceInterface) *VenueHandler {
	return &VenueHandler{venueService: s}
}

func (h *VenueHandler) HandleCreateVenue(c *gin.Context) {
	var req services.CreateVenueRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newID, err := h.venueService.CreateVenue(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

func (h *VenueHandler) HandleGetVenueByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	venue, err := h.venueService.GetVenueByID(c.Request.Context(), id)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, toDTOVenue(*venue))
}

func (h *VenueHandler) HandleListVenues(c *gin.Context) {
	venues, err := h.venueService.ListVenues(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	venueDTOs := make([]venueDTO, 0, len(venues))
	for _, v := range venues {
		venueDTOs = append(venueDTOs, toDTOVenue(v))
	}
	c.JSON(http.StatusOK, venueDTOs)
}

func (h *VenueHandler) HandleUpdateVenue(c *gin.Context) {
	var req services.UpdateVenueRequest

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = h.venueService.UpdateVenue(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *VenueHandler) HandleDeleteVenue(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	err = h.venueService.DeleteVenue(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
