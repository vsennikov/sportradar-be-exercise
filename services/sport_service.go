package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type SportRepositoryInterface interface {
	CreateSport(ctx context.Context, name string) (int, error)
	GetSportById(ctx context.Context, id int) (*Sport, error)
	ListSports(ctx context.Context) ([]Sport, error)
	UpdateSport(ctx context.Context, id int, name string) error
	DeleteSport(ctx context.Context, id int) error
}

type SportServiceInterface interface {
	CreateSport(ctx context.Context, req SportRequest) (int, error)
	GetSportByID(ctx context.Context, id int) (*Sport, error)
	ListSports(ctx context.Context) ([]Sport, error)
	UpdateSport(ctx context.Context, id int, req SportRequest) error
	DeleteSport(ctx context.Context, id int) error
}

type SportService struct {
	sportRepository SportRepositoryInterface
	eventRepository EventRepositoryInterface
}

func NewSportService(r SportRepositoryInterface, e EventRepositoryInterface) *SportService{
	return &SportService{sportRepository: r, eventRepository: e}
}

func (s *SportService) CreateSport(ctx context.Context, req SportRequest) (int, error) {
	if len(req.Name) < 3 {
		return 0, fmt.Errorf("sport name must be at least 3 characters long")
	}
	newID, err := s.sportRepository.CreateSport(ctx, req.Name)
	if err != nil {
		return 0, fmt.Errorf("failed to create sport: %w", err)
	}
	return newID, nil
}

func (s *SportService) GetSportByID(ctx context.Context, id int) (*Sport, error) {
	sport, err := s.sportRepository.GetSportById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("sport with id %d not found", id)
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	return sport, nil
}

func (s *SportService) ListSports(ctx context.Context) ([]Sport, error) {
	sports, err := s.sportRepository.ListSports(ctx)
	if err != nil {
	return nil, fmt.Errorf("failed to list sports: %w", err)
	}
	return sports, nil
}

func (s *SportService) UpdateSport(ctx context.Context, id int, req SportRequest) error {
	if len(req.Name) < 3 {
		return fmt.Errorf("sport name must be at least 3 characters long")
	}
	sportToUpdate := Sport{
		ID: id,
		Name: req.Name,
	}
	err := s.sportRepository.UpdateSport(ctx, sportToUpdate.ID, sportToUpdate.Name)
	if err != nil {
		return fmt.Errorf("failed to update sport: %w", err)
	}
	return nil
}

func (s *SportService) DeleteSport(ctx context.Context, id int) error {
	count, err := s.eventRepository.CountEventsBySportID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check event usage: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete sport: it is currently used by %d events", count)
	}
	err = s.sportRepository.DeleteSport(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete sport: %w", err)
	}
	return nil
}