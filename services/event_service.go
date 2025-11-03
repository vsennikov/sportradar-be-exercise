package services

import (
	"context"
	"fmt"
	"math"
	"time"
)

type EventRepositoryInterface interface {
	GetEventByID(ctx context.Context, id int) (*Event, error)
	CreateEvent(ctx context.Context, params CreateEventParams) (int, error)
	CountEvents(ctx context.Context, params ListEventsParams) (int, error)
	ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error)
	CountEventsBySportID(ctx context.Context, sportID int) (int, error)
	CountEventsByVenueId(ctx context.Context, venueID int) (int, error)
	CountEventsByTeamID(ctx context.Context, teamID int) (int, error)
	UpdateEvent(ctx context.Context, event Event) error
	DeleteEvent(ctx context.Context, id int) error
}

type EventServiceInterface interface {
	GetEventByID(ctx context.Context, id int) (*Event, error)
	CreateEvent(ctx context.Context, req EventCreateRequest) (int, error)
	ListEvents(ctx context.Context, req ListEventsRequest) ([]Event, *Pagination, error)
	UpdateEvent(ctx context.Context, id int, req UpdateEventRequest) error
	DeleteEvent(ctx context.Context, id int) error
}

type EventService struct {
	eventRepository   EventRepositoryInterface
	defaultPage  int
	defaultLimit int
	sportRepository SportRepositoryInterface
	teamRepository TeamRepositoryInterface
	venueRepository VenueRepositoryInterface
}

func NewEventService(r EventRepositoryInterface, dP, dL int,
	 s SportRepositoryInterface, t TeamRepositoryInterface, v VenueRepositoryInterface) *EventService {
	return &EventService{
		eventRepository: r,
		defaultPage: dP,
		defaultLimit: dL,
		sportRepository: s,
		teamRepository: t,
		venueRepository: v,}
}

func (s *EventService) GetEventByID(ctx context.Context, id int) (*Event, error) {
	event, err := s.eventRepository.GetEventByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("database error while fetching event: %w", err)
	}
	return event, nil
}

func (s *EventService) CreateEvent(ctx context.Context, req EventCreateRequest) (int, error) {
	if req.EventDatetime.Before(time.Now()) {
		return 0, fmt.Errorf("cannot create an event in the past")
	}
	params := CreateEventParams(req)
	newID, err := s.eventRepository.CreateEvent(ctx, params)
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
	totalItems, err := s.eventRepository.CountEvents(ctx, repoParams)
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
	events, err := s.eventRepository.ListEvents(ctx, repoParams)
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

func (s *EventService) UpdateEvent(ctx context.Context, id int, req UpdateEventRequest) error {
	existingEvent, err := s.eventRepository.GetEventByID(ctx, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if req.EventDatetime != nil {
		existingEvent.EventDatetime = *req.EventDatetime
	}
	if req.Description != nil {
		existingEvent.Description = req.Description
	}
	if req.HomeScore != nil {
		existingEvent.HomeScore = req.HomeScore
	}
	if req.AwayScore != nil {
		existingEvent.AwayScore = req.AwayScore
	}
	if req.SportID != nil {
		sport, err := s.sportRepository.GetSportById(ctx, *req.SportID)
		if err != nil {
			return fmt.Errorf("validation error: sport with id %d not found", *req.SportID)
		}
		existingEvent.Sport = *sport
	}
	if req.VenueID != nil {
		venue, err := s.venueRepository.GetVenueById(ctx, *req.VenueID)
		if err != nil {
			return fmt.Errorf("validation error: venue with id %d not found", *req.VenueID)
		}
		existingEvent.Venue = *venue
	}
	if req.HomeTeamID != nil {
		team, err := s.teamRepository.GetTeamByID(ctx, *req.HomeTeamID)
		if err != nil {
			return fmt.Errorf("validation error: home team with id %d not found", *req.HomeTeamID)
		}
		existingEvent.HomeTeam = *team
	}
	if req.AwayTeamID != nil {
				team, err := s.teamRepository.GetTeamByID(ctx, *req.AwayTeamID)
		if err != nil {
			return fmt.Errorf("validation error: away team with id %d not found", *req.AwayTeamID)
		}
		existingEvent.AwayTeam = *team
	}
	err = s.eventRepository.UpdateEvent(ctx, *existingEvent)
	if err != nil {
		return fmt.Errorf("failed to update event: %w", err)
	}
	return nil
}

func (s *EventService) DeleteEvent(ctx context.Context, id int) error {
	_, err := s.eventRepository.GetEventByID(ctx, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	err = s.eventRepository.DeleteEvent(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete event: %w", err)
	}
	return nil
}
