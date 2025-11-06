package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockTeamRepositoryForService is a mock for TeamRepositoryInterface
type MockTeamRepositoryForService struct {
	mock.Mock
}

func (m *MockTeamRepositoryForService) CreateTeam(ctx context.Context, params TeamRequest) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockTeamRepositoryForService) GetTeamByID(ctx context.Context, id int) (*Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Team), args.Error(1)
}

func (m *MockTeamRepositoryForService) ListTeams(ctx context.Context) ([]Team, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Team), args.Error(1)
}

func (m *MockTeamRepositoryForService) UpdateTeam(ctx context.Context, team Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamRepositoryForService) DeleteTeam(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEventRepositoryForTeam is a mock for EventRepositoryInterface used in TeamService tests
type MockEventRepositoryForTeam struct {
	mock.Mock
}

func (m *MockEventRepositoryForTeam) GetEventByID(ctx context.Context, id int) (*Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Event), args.Error(1)
}

func (m *MockEventRepositoryForTeam) CreateEvent(ctx context.Context, params CreateEventParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForTeam) CountEvents(ctx context.Context, params ListEventsParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForTeam) ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Event), args.Error(1)
}

func (m *MockEventRepositoryForTeam) CountEventsBySportID(ctx context.Context, sportID int) (int, error) {
	args := m.Called(ctx, sportID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForTeam) CountEventsByVenueId(ctx context.Context, venueID int) (int, error) {
	args := m.Called(ctx, venueID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForTeam) CountEventsByTeamID(ctx context.Context, teamID int) (int, error) {
	args := m.Called(ctx, teamID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForTeam) UpdateEvent(ctx context.Context, event Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepositoryForTeam) DeleteEvent(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestTeamService_CreateTeam(t *testing.T) {
	tests := []struct {
		name          string
		request       CreateTeamRequest
		mockID        int
		mockError     error
		expectedID    int
		expectedError bool
	}{
		{
			name:          "successful creation",
			request:      CreateTeamRequest{Name: "Lakers", City: "Los Angeles", SportID: 1},
			mockID:        1,
			mockError:     nil,
			expectedID:    1,
			expectedError: false,
		},
		{
			name:          "name too short",
			request:      CreateTeamRequest{Name: "AB", City: "Los Angeles", SportID: 1},
			mockID:        0,
			mockError:     nil,
			expectedID:    0,
			expectedError: true,
		},
		{
			name:          "city too short",
			request:      CreateTeamRequest{Name: "Lakers", City: "LA", SportID: 1},
			mockID:        0,
			mockError:     nil,
			expectedID:    0,
			expectedError: true,
		},
		{
			name:          "database error",
			request:      CreateTeamRequest{Name: "Lakers", City: "Los Angeles", SportID: 1},
			mockID:        0,
			mockError:     errors.New("database error"),
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTeamRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForTeam)

			service := NewTeamService(mockRepo, mockEventRepo)

			if !tt.expectedError || tt.name == "database error" {
				mockRepo.On("CreateTeam", mock.Anything, mock.AnythingOfType("TeamRequest")).Return(tt.mockID, tt.mockError)
			}

			result, err := service.CreateTeam(context.Background(), tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Equal(t, 0, result)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedID, result)
			}

			if !tt.expectedError || tt.name == "database error" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestTeamService_GetTeamByID(t *testing.T) {
	tests := []struct {
		name          string
		teamID        int
		mockTeam      *Team
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful retrieval",
			teamID:        1,
			mockTeam:      &Team{ID: 1, Name: "Lakers", City: "Los Angeles"},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "team not found",
			teamID:        999,
			mockTeam:      nil,
			mockError:     sql.ErrNoRows,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTeamRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForTeam)

			service := NewTeamService(mockRepo, mockEventRepo)

			mockRepo.On("GetTeamByID", mock.Anything, tt.teamID).Return(tt.mockTeam, tt.mockError)

			result, err := service.GetTeamByID(context.Background(), tt.teamID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.teamID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTeamService_ListTeams(t *testing.T) {
	tests := []struct {
		name          string
		mockTeams     []Team
		mockError     error
		expectedCount int
		expectedError bool
	}{
		{
			name:          "successful list",
			mockTeams:     []Team{{ID: 1, Name: "Lakers"}, {ID: 2, Name: "Warriors"}},
			mockError:     nil,
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "empty list",
			mockTeams:     []Team{},
			mockError:     nil,
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "database error",
			mockTeams:     nil,
			mockError:     errors.New("database error"),
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTeamRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForTeam)

			service := NewTeamService(mockRepo, mockEventRepo)

			mockRepo.On("ListTeams", mock.Anything).Return(tt.mockTeams, tt.mockError)

			result, err := service.ListTeams(context.Background())

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.expectedCount, len(result))
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTeamService_UpdateTeam(t *testing.T) {
	existingTeam := &Team{
		ID:   1,
		Name: "Lakers",
		City: "Los Angeles",
	}

	tests := []struct {
		name          string
		teamID        int
		request       UpdateTeamRequest
		mockTeam      *Team
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful update",
			teamID:        1,
			request:       UpdateTeamRequest{Name: stringPtr("Updated Lakers")},
			mockTeam:      existingTeam,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "name too short",
			teamID:        1,
			request:       UpdateTeamRequest{Name: stringPtr("AB")},
			mockTeam:      existingTeam,
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "city too short",
			teamID:        1,
			request:       UpdateTeamRequest{City: stringPtr("LA")},
			mockTeam:      existingTeam,
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "team not found",
			teamID:        999,
			request:       UpdateTeamRequest{Name: stringPtr("Updated")},
			mockTeam:      nil,
			mockError:     sql.ErrNoRows,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTeamRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForTeam)

			service := NewTeamService(mockRepo, mockEventRepo)

			mockRepo.On("GetTeamByID", mock.Anything, tt.teamID).Return(tt.mockTeam, tt.mockError)

			if tt.mockTeam != nil && !tt.expectedError {
				mockRepo.On("UpdateTeam", mock.Anything, mock.AnythingOfType("Team")).Return(nil)
			}

			err := service.UpdateTeam(context.Background(), tt.teamID, tt.request)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestTeamService_DeleteTeam(t *testing.T) {
	tests := []struct {
		name          string
		teamID        int
		eventCount    int
		countError    error
		deleteError   error
		expectedError bool
	}{
		{
			name:          "successful deletion",
			teamID:        1,
			eventCount:    0,
			countError:    nil,
			deleteError:   nil,
			expectedError: false,
		},
		{
			name:          "team in use",
			teamID:        1,
			eventCount:    5,
			countError:    nil,
			deleteError:   nil,
			expectedError: true,
		},
		{
			name:          "count error",
			teamID:        1,
			eventCount:    0,
			countError:    errors.New("database error"),
			deleteError:   nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockTeamRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForTeam)

			service := NewTeamService(mockRepo, mockEventRepo)

			mockEventRepo.On("CountEventsByTeamID", mock.Anything, tt.teamID).Return(tt.eventCount, tt.countError)

			if tt.countError == nil && tt.eventCount == 0 {
				mockRepo.On("DeleteTeam", mock.Anything, tt.teamID).Return(tt.deleteError)
			}

			err := service.DeleteTeam(context.Background(), tt.teamID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockEventRepo.AssertExpectations(t)
			if tt.countError == nil && tt.eventCount == 0 {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

