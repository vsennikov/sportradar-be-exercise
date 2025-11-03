package services

import (
	"context"
	"fmt"
)

type VenueRepositoryInterface interface {
	CreateVenue(ctx context.Context, params VenueRequest) (int, error)
	GetVenueById(ctx context.Context, id int) (*Venue, error)
	ListVenues(ctx context.Context) ([]Venue, error)
	UpdateVenue(ctx context.Context, venue Venue) error
	DeleteVenue(ctx context.Context, id int) error
}

type VenueServiceInterface interface {
	CreateVenue(ctx context.Context, req CreateVenueRequest) (int, error)
	GetVenueByID(ctx context.Context, id int) (*Venue, error)
	ListVenues(ctx context.Context) ([]Venue, error)
	UpdateVenue(ctx context.Context, id int, req UpdateVenueRequest) error
	DeleteVenue(ctx context.Context, id int) error
}

type VenueService struct {
	venueRepository VenueRepositoryInterface
	eventRepository EventRepositoryInterface
}

func NewVenueService(v VenueRepositoryInterface, e EventRepositoryInterface) *VenueService{
	return &VenueService{venueRepository: v, eventRepository: e}
}

func (s *VenueService) CreateVenue(ctx context.Context, req CreateVenueRequest) (int, error) {
	if len(req.Name) < 3 {
		return 0, fmt.Errorf("venue name must be at least 3 characters long")
	} else if len (req.City) < 3 {
		return 0, fmt.Errorf("venue city name must be at least 3 characters long")

	} else if len (req.CountryCode) != 2 {
		return 0, fmt.Errorf("venue city_code must be 2 characters long")
	}
	params := VenueRequest(req)
	newID, err := s.venueRepository.CreateVenue(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create venue: %w", err)
	}
	return newID, nil
}

func (s *VenueService) GetVenueByID(ctx context.Context, id int) (*Venue, error) {
	venue, err := s.venueRepository.GetVenueById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return venue, nil
}

func (s *VenueService) ListVenues(ctx context.Context) ([]Venue, error) {
	venues, err := s.venueRepository.ListVenues(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list venues: %w", err)
	}
	return venues, nil
}

func (s *VenueService) UpdateVenue(ctx context.Context, id int, req UpdateVenueRequest) error {
	existingVenue, err := s.venueRepository.GetVenueById(ctx, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if req.Name != nil {
		existingVenue.Name = *req.Name
	}
	if req.City != nil {
		existingVenue.City = *req.City
	}
	if req.CountryCode != nil {
		existingVenue.CountryCode = *req.CountryCode
	}
	err = s.venueRepository.UpdateVenue(ctx, *existingVenue)
	if err != nil {
		return fmt.Errorf("failed to update venue: %w", err)
	}
	return nil
}

func (s *VenueService) DeleteVenue(ctx context.Context, id int) error {
	count, err := s.eventRepository.CountEventsByVenueId(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check event usage: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete venue: it is currently used by %d events", count)
	}
	err = s.venueRepository.DeleteVenue(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete venue: %w", err)
	}
	return nil
}
