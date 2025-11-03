package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type SportHandler struct {
	sportService services.SportServiceInterface
}

func NewSportHandler(s services.SportServiceInterface) *SportHandler {
	return &SportHandler{sportService: s}
}

func (h *SportHandler) HandleCreateSport(c *gin.Context) {
	var req services.SportRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newID, err := h.sportService.CreateSport(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

func (h *SportHandler) HandleGetSportByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	sport, err := h.sportService.GetSportByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, toDTOSport(*sport))
}

func (h *SportHandler) HandleListSports(c *gin.Context) {
	sports, err := h.sportService.ListSports(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sportsDTO := make([]sportDTO, 0, len(sports))
	for _, sport := range sports {
		sportsDTO = append(sportsDTO,toDTOSport(sport))
	}

	c.JSON(http.StatusOK, sportsDTO)
}

func (h *SportHandler) HandleUpdateSport(c *gin.Context) {
	var req services.SportRequest

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
	err = h.sportService.UpdateSport(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *SportHandler) HandleDeleteSport(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}

	err = h.sportService.DeleteSport(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}