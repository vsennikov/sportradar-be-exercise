package infrastructure

import (
	"context"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type EventRepository struct {
	db *sqlx.DB
}

func NewEventRepository(db *sqlx.DB) *EventRepository {
	return &EventRepository{
		db: db,
	}
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

func (r *EventRepository) CreateEvent(ctx context.Context, params services.CreateEventParams) (int, error) {
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

func (r *EventRepository) CountEvents(ctx context.Context, params services.ListEventsParams) (int, error) {
	var args []interface{}
	var total int
	i := 1
	whereQuery := []string{"1=1"}

	if params.SportID != nil {
		args = append(args, *params.SportID)
		whereQuery = append(whereQuery, fmt.Sprintf("e._sport_id = $%d", i))
		i++
	}
	if params.DateFrom != nil {
		args = append(args, *params.DateFrom)
		whereQuery = append(whereQuery, fmt.Sprintf("event_datetime >= $%d", i))
		i++
	}
	query := fmt.Sprintf("SELECT COUNT(*) FROM events e WHERE %s", strings.Join(whereQuery, " AND "))
	if err := r.db.GetContext(ctx, &total, query, args...); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *EventRepository) ListEvents(ctx context.Context,
	params services.ListEventsParams) ([]services.Event, error) {
	var dbModels []eventDBModel
	var args []interface{}
	i := 1
	whereQuery := []string{"1=1"}

	if params.SportID != nil {
		args = append(args, *params.SportID)
		whereQuery = append(whereQuery, fmt.Sprintf("e._sport_id = $%d", i))
		i++
	}
	if params.DateFrom != nil {
		args = append(args, *params.DateFrom)
		whereQuery = append(whereQuery, fmt.Sprintf("e.event_datetime >= $%d", i))
		i++
	}
	args = append(args, params.Limit)
	limitClause := fmt.Sprintf("LIMIT $%d", i)
	i++
	args = append(args, params.Offset)
	offsetClause := fmt.Sprintf("OFFSET $%d", i)
	i++
	query := fmt.Sprintf(
		"%s WHERE %s ORDER BY e.event_datetime ASC %s %s",
		baseEventSelectQuery,
		strings.Join(whereQuery, " AND "),
		limitClause,
		offsetClause,
	)
	if err := r.db.SelectContext(ctx, &dbModels, query, args...); err != nil {
		return nil, err
	}
	events := make([]services.Event, 0, len(dbModels))
	for _, dbModel := range dbModels {
		events = append(events, toServiceEvent(dbModel))
	}
	return events, nil
}

func (r *EventRepository) CountEventsBySportID(ctx context.Context, sportID int) (int, error) {
	query := "SELECT COUNT(*) FROM events WHERE _sport_id = $1"
	var total int

	if err := r.db.GetContext(ctx, &total, query, sportID); err != nil {
		return 0, err
	}
	return total, nil
}

func (r *EventRepository) CountEventsByVenueId(ctx context.Context, venueID int) (int, error) {
	query := "SELECT COUNT(*) FROM events WHERE _venue_id = $1"
	var total int

	if err := r.db.GetContext(ctx, &total, query, venueID); err != nil {
		return 0, err
	}
	return total, nil
}
