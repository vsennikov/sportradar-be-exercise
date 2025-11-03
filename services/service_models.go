package services

import "time"

type Sport struct {
	ID   int
	Name string
}

type Venue struct {
	ID          int
	Name        string
	City        string
	CountryCode string
}

type Team struct {
	ID   int
	Name string
	City string
}

type Event struct {
	ID            int
	EventDatetime time.Time
	Description   *string
	HomeScore     *int
	AwayScore     *int

	Sport    Sport
	Venue    Venue
	HomeTeam Team
	AwayTeam Team
}

type ListEventsParams struct {
	SportID  *int
	DateFrom *time.Time
	Limit    int
	Offset   int
}

type CreateEventParams struct {
	EventDatetime time.Time
	Description   *string
	SportID       int
	VenueID       *int
	HomeTeamID    int
	AwayTeamID    int
}

type ListEventsRequest struct {
	SportID  *int
	DateFrom *time.Time
	Page     int
	Limit    int
}

type Pagination struct {
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

type EventCreateRequest struct {
	EventDatetime time.Time `json:"event_datetime" binding:"required"`
	Description   *string   `json:"description"`
	SportID       int       `json:"sport_id" binding:"required"`
	VenueID       *int      `json:"venue_id"`
	HomeTeamID    int       `json:"home_team_id" binding:"required"`
	AwayTeamID    int       `json:"away_team_id" binding:"required"`
}

type SportRequest struct {
	Name string `json:"name" binding:"required"`
}

type VenueRequest struct {
	Name string
	City string
	CountryCode string
}

type CreateVenueRequest struct {
	Name        string `json:"name" binding:"required"`
	City        string `json:"city" binding:"required"`
	CountryCode string `json:"country_code" binding:"required"`
}

type UpdateVenueRequest struct {
	Name        *string `json:"name"`
	City        *string `json:"city"`
	CountryCode *string `json:"country_code"`
}