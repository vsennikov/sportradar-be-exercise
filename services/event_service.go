package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"time"
)

type EventRepositoryInterface interface {
	GetEventByID(ctx context.Context, id int) (*Event, error)
	CreateEvent(ctx context.Context, params CreateEventParams) (int, error)
	CountEvents(ctx context.Context, params ListEventsParams) (int, error)
	ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error)
}

type EventServiceInterface interface {
	GetEventByID(ctx context.Context, id int) (*Event, error)
	CreateEvent(ctx context.Context, req EventCreateRequest) (int, error)
	ListEvents(ctx context.Context, req ListEventsRequest) ([]Event, *Pagination, error)
}

type EventService struct {
	repository EventRepositoryInterface
	defaultPage int
	defaultLimit int
}

func NewEventService(repository EventRepositoryInterface, defaultPage, defaultLimit int) *EventService {
    return &EventService{repository: repository, defaultPage: defaultPage, defaultLimit: defaultLimit}
}

func (s *EventService) GetEventByID(ctx context.Context, id int) (*Event, error) {
	event, err := s.repository.GetEventByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("event with id %d not found", id)
		}
		return nil, fmt.Errorf("database error while fetching event: %w", err)
	}
	return event, nil
}

func (s *EventService) CreateEvent(ctx context.Context, req EventCreateRequest) (int, error) {
	if req.EventDatetime.Before(time.Now()) {
		return 0, fmt.Errorf("cannot create an event in the past")
	}
	params := CreateEventParams{
		EventDatetime: req.EventDatetime,
		Description:   req.Description,
		SportID:       req.SportID,
		VenueID:       req.VenueID,
		HomeTeamID:    req.HomeTeamID,
		AwayTeamID:    req.AwayTeamID,
	}
	newID, err := s.repository.CreateEvent(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create event: %w", err)
	}
	return newID, nil
}

func (s *EventService) ListEvents(ctx context.Context, req ListEventsRequest) ([]Event, *Pagination, error) {
	if req.Page <= 0 {
		req.Page = s.defaultPage
	}
	if req.Limit <= 0 {
		req.Limit = s.defaultLimit
	}
	offset := (req.Page - 1) * req.Limit
	repoParams := ListEventsParams{
		SportID:  req.SportID,
		DateFrom: req.DateFrom,
		Limit:    req.Limit,
		Offset:   offset,
	}
	totalItems, err := s.repository.CountEvents(ctx, repoParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count events: %w", err)
	}
	if totalItems == 0 {
		pagination := &Pagination{
			TotalItems:  0,
			TotalPages:  0,
			CurrentPage: req.Page,
			PageSize:    req.Limit,
		}
		return []Event{}, pagination, nil
	}
	events, err := s.repository.ListEvents(ctx, repoParams)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list events: %w", err)
	}
	totalPages := int(math.Ceil(float64(totalItems) / float64(req.Limit)))
	pagination := &Pagination{
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: req.Page,
		PageSize:    req.Limit,
	}
	return events, pagination, nil
}