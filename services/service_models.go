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
