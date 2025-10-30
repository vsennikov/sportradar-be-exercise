package controllers

import "time"

type sportDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type venueDTO struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	City        string `json:"city"`
	CountryCode string `json:"country_code"`
}

type teamDTO struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	City string `json:"city"`
}

type EventDTO struct {
	ID            int        `json:"id"`
	EventDatetime time.Time  `json:"event_datetime"`
	Description   *string    `json:"description,omitempty"`
	HomeScore     *int       `json:"home_score,omitempty"`
	AwayScore     *int       `json:"away_score,omitempty"`
	Sport         sportDTO   `json:"sport"`
	Venue         *venueDTO  `json:"venue,omitempty"`
	HomeTeam      teamDTO    `json:"home_team"`
	AwayTeam      teamDTO    `json:"away_team"`
}

type EventCreateRequest struct {
	EventDatetime time.Time `json:"event_datetime" binding:"required"`
	Description   *string   `json:"description"`
	SportID       int       `json:"sport_id" binding:"required"`
	VenueID       *int      `json:"venue_id"`
	HomeTeamID    int       `json:"home_team_id" binding:"required"`
	AwayTeamID    int       `json:"away_team_id" binding:"required"`
}