package infrastructure

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
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

func (r *EventRepository) GetEventByID(ctx context.Context, id int) (*services.Event, error) {
	var dbModel eventDBModel
	query := baseEventSelectQuery + " WHERE e.id = $1"

	if err := r.db.GetContext(ctx, &dbModel, query, id); err != nil {
		return nil, err
	}

	event := toServiceEvent(dbModel)
	return &event, nil
}

func (r *EventRepository) CreateEvent(ctx context.Context, params CreateEventParams) (int, error) {
	var newID int

	query := `
	INSERT INTO events(event_datetime, description, _sport_id, _venue_id, _home_team_id, _away_team_id)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id`

	err := r.db.QueryRowContext(
		ctx, query,
		params.EventDatetime, params.Description, params.SportID,
		params.VenueID, params.HomeTeamID, params.AwayTeamID,
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

