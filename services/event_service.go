package services

import (
	"context"
)

type EventRepository interface {
	GetEventByID(ctx context.Context, id int) (*Event, error)
	CreateEvent(ctx context.Context, params CreateEventParams) (int, error)
	CountEvents(ctx context.Context, params ListEventsParams) (int, error)
	ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error)
}



type EventService struct {
	repository EventRepository
}

func NewEventService(repository EventRepository) *EventService {
    return &EventService{repository: repository}
}