package infrastructure

import (
	"database/sql"
	"time"
)


type eventDBModel struct {
	ID            int            `db:"id"`
	EventDatetime time.Time      `db:"event_datetime"`
	Description   sql.NullString `db:"description"`
	HomeScore     sql.NullInt64  `db:"home_score"`
	AwayScore     sql.NullInt64  `db:"away_score"`

	SportID   int    `db:"sport.id"`
	SportName string `db:"sport.name"`

	VenueID          sql.NullInt64  `db:"venue.id"`
	VenueName        sql.NullString `db:"venue.name"`
	VenueCity        sql.NullString `db:"venue.city"`
	VenueCountryCode sql.NullString `db:"venue.country_code"`

	HomeTeamID   int    `db:"ht.id"`
	HomeTeamName string `db:"ht.name"`
	HomeTeamCity string `db:"ht.city"`

	AwayTeamID   int    `db:"at.id"`
	AwayTeamName string `db:"at.name"`
	AwayTeamCity string `db:"at.city"`
}

type sportDBModel struct {
	ID	 int	`db:"id"`
	Name string `db:"name"`
}