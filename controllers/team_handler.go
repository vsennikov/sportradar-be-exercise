package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type TeamHandler struct {
	teamService services.TeamServiceInterface
}

func NewTeamHandler(s services.TeamServiceInterface) *TeamHandler {
	return &TeamHandler{teamService: s}
}

func (h *TeamHandler) HandleCreateTeam(c *gin.Context) {
	var req services.CreateTeamRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newID, err := h.teamService.CreateTeam(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID})
}

func (h *TeamHandler) HandleGetTeamByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	team, err := h.teamService.GetTeamByID(c.Request.Context(), id)
	if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
	}
	c.JSON(http.StatusOK,toDTOTeam(*team))
}

func (h *TeamHandler) HandleListTeams(c *gin.Context) {
	teams, err := h.teamService.ListTeams(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	teamDTOs := make([]teamDTO, 0, len(teams))
	for _, t := range teams {
		teamDTOs = append(teamDTOs, toDTOTeam(t))
	}
	c.JSON(http.StatusOK, teamDTOs)
}

func (h *TeamHandler) HandleUpdateTeam(c *gin.Context) {
	var req services.UpdateTeamRequest

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
	err = h.teamService.UpdateTeam(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *TeamHandler) HandleDeleteTeam(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID format"})
		return
	}
	err = h.teamService.DeleteTeam(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}
