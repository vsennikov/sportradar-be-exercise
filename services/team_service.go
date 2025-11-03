package services

import (
	"context"
	"fmt"
)

type TeamRepositoryInterface interface{
	CreateTeam(ctx context.Context, params TeamRequest) (int, error)
	GetTeamByID(ctx context.Context, id int) (*Team, error)
	ListTeams(ctx context.Context) ([]Team, error)
	UpdateTeam(ctx context.Context, team Team) error
	DeleteTeam(ctx context.Context, id int) error
}

type TeamServiceInterface interface{
	CreateTeam(ctx context.Context, req CreateTeamRequest) (int, error)
	GetTeamByID(ctx context.Context, id int) (*Team, error)
	ListTeams(ctx context.Context) ([]Team, error)
	UpdateTeam(ctx context.Context, id int, req UpdateTeamRequest) error
	DeleteTeam(ctx context.Context, id int) error
}

type TeamService struct{
	teamRepository TeamRepositoryInterface
	eventRepository EventRepositoryInterface
}

func NewTeamService (t TeamRepositoryInterface, e EventRepositoryInterface) *TeamService{
	return &TeamService{teamRepository: t, eventRepository: e}
}

func (s *TeamService) CreateTeam(ctx context.Context, req CreateTeamRequest) (int, error) {
	if len(req.Name) < 3 {
		return 0, fmt.Errorf("team name must be at least 3 characters long")
	} else if len (req.City) < 3 {
		return 0, fmt.Errorf("team city name must be at least 3 characters long")
	}
	params := TeamRequest(req)
	newID, err := s.teamRepository.CreateTeam(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to create team: %w", err)
	}
	return newID, nil
}

func (s *TeamService) GetTeamByID(ctx context.Context, id int) (*Team, error) {
	team, err := s.teamRepository.GetTeamByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}
	return team, nil
}

func (s *TeamService) ListTeams(ctx context.Context) ([]Team, error) {
	teams, err := s.teamRepository.ListTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	return teams, nil
}

func (s *TeamService) UpdateTeam(ctx context.Context, id int, req UpdateTeamRequest) error {
	existingTeam, err := s.teamRepository.GetTeamByID(ctx, id)
	if err != nil {
		return fmt.Errorf("database error: %w", err)
	}
	if req.Name != nil {
		if len(*req.Name) < 3 {
			return fmt.Errorf("team name must be at least 3 characters long")
		}
		existingTeam.Name = *req.Name
	}
	if req.City != nil {
		if len(*req.City) < 3 {
			return fmt.Errorf("team city name must be at least 3 characters long")
		}
		existingTeam.City = *req.City
	}
	err = s.teamRepository.UpdateTeam(ctx, *existingTeam)
	if err != nil {
		return fmt.Errorf("failed to update team: %w", err)
	}
	return nil
}

func (s *TeamService) DeleteTeam(ctx context.Context, id int) error {
	count, err := s.eventRepository.CountEventsByTeamID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to check event usage: %w", err)
	}
	if count > 0 {
		return fmt.Errorf("cannot delete team: it is currently used by %d events", count)
	}
	err = s.teamRepository.DeleteTeam(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete team: %w", err)
	}
	return nil
}
