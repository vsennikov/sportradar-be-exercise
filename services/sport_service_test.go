package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepositoryForSport is a mock for EventRepositoryInterface used in SportService tests
type MockEventRepositoryForSport struct {
	mock.Mock
}

func (m *MockEventRepositoryForSport) GetEventByID(ctx context.Context, id int) (*Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Event), args.Error(1)
}

func (m *MockEventRepositoryForSport) CreateEvent(ctx context.Context, params CreateEventParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForSport) CountEvents(ctx context.Context, params ListEventsParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForSport) ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Event), args.Error(1)
}

func (m *MockEventRepositoryForSport) CountEventsBySportID(ctx context.Context, sportID int) (int, error) {
	args := m.Called(ctx, sportID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForSport) CountEventsByVenueId(ctx context.Context, venueID int) (int, error) {
	args := m.Called(ctx, venueID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForSport) CountEventsByTeamID(ctx context.Context, teamID int) (int, error) {
	args := m.Called(ctx, teamID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForSport) UpdateEvent(ctx context.Context, event Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepositoryForSport) DeleteEvent(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSportRepositoryForService is a mock for SportRepositoryInterface
type MockSportRepositoryForService struct {
	mock.Mock
}

func (m *MockSportRepositoryForService) CreateSport(ctx context.Context, name string) (int, error) {
	args := m.Called(ctx, name)
	return args.Int(0), args.Error(1)
}

func (m *MockSportRepositoryForService) GetSportById(ctx context.Context, id int) (*Sport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Sport), args.Error(1)
}

func (m *MockSportRepositoryForService) ListSports(ctx context.Context) ([]Sport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Sport), args.Error(1)
}

func (m *MockSportRepositoryForService) UpdateSport(ctx context.Context, id int, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockSportRepositoryForService) DeleteSport(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSportService_CreateSport(t *testing.T) {
	tests := []struct {
		name          string
		request       SportRequest
		mockID        int
		mockError     error
		expectedID    int
		expectedError bool
	}{
		{
			name:          "successful creation",
			request:      SportRequest{Name: "Basketball"},
			mockID:        1,
			mockError:     nil,
			expectedID:    1,
			expectedError: false,
		},
		{
			name:          "name too short",
			request:      SportRequest{Name: "AB"},
			mockID:        0,
			mockError:     nil,
			expectedID:    0,
			expectedError: true,
		},
		{
			name:          "database error",
			request:      SportRequest{Name: "Basketball"},
			mockID:        0,
			mockError:     errors.New("database error"),
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSportRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForSport)

			service := NewSportService(mockRepo, mockEventRepo)

			if !tt.expectedError || tt.name == "database error" {
				mockRepo.On("CreateSport", mock.Anything, tt.request.Name).Return(tt.mockID, tt.mockError)
			}

			result, err := service.CreateSport(context.Background(), tt.request)

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

func TestSportService_GetSportByID(t *testing.T) {
	tests := []struct {
		name          string
		sportID       int
		mockSport     *Sport
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful retrieval",
			sportID:       1,
			mockSport:     &Sport{ID: 1, Name: "Football"},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "sport not found",
			sportID:       999,
			mockSport:     nil,
			mockError:     sql.ErrNoRows,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSportRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForSport)

			service := NewSportService(mockRepo, mockEventRepo)

			mockRepo.On("GetSportById", mock.Anything, tt.sportID).Return(tt.mockSport, tt.mockError)

			result, err := service.GetSportByID(context.Background(), tt.sportID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.sportID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSportService_ListSports(t *testing.T) {
	tests := []struct {
		name          string
		mockSports    []Sport
		mockError     error
		expectedCount int
		expectedError bool
	}{
		{
			name:          "successful list",
			mockSports:    []Sport{{ID: 1, Name: "Football"}, {ID: 2, Name: "Basketball"}},
			mockError:     nil,
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "empty list",
			mockSports:    []Sport{},
			mockError:     nil,
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "database error",
			mockSports:    nil,
			mockError:     errors.New("database error"),
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSportRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForSport)

			service := NewSportService(mockRepo, mockEventRepo)

			mockRepo.On("ListSports", mock.Anything).Return(tt.mockSports, tt.mockError)

			result, err := service.ListSports(context.Background())

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

func TestSportService_UpdateSport(t *testing.T) {
	tests := []struct {
		name          string
		sportID       int
		request       SportRequest
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful update",
			sportID:       1,
			request:       SportRequest{Name: "Updated Football"},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "name too short",
			sportID:       1,
			request:       SportRequest{Name: "AB"},
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "database error",
			sportID:       1,
			request:       SportRequest{Name: "Updated Football"},
			mockError:     errors.New("database error"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSportRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForSport)

			service := NewSportService(mockRepo, mockEventRepo)

			if !tt.expectedError || tt.name == "database error" {
				mockRepo.On("UpdateSport", mock.Anything, tt.sportID, tt.request.Name).Return(tt.mockError)
			}

			err := service.UpdateSport(context.Background(), tt.sportID, tt.request)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if !tt.expectedError || tt.name == "database error" {
				mockRepo.AssertExpectations(t)
			}
		})
	}
}

func TestSportService_DeleteSport(t *testing.T) {
	tests := []struct {
		name          string
		sportID       int
		eventCount    int
		countError    error
		deleteError   error
		expectedError bool
	}{
		{
			name:          "successful deletion",
			sportID:       1,
			eventCount:    0,
			countError:    nil,
			deleteError:   nil,
			expectedError: false,
		},
		{
			name:          "sport in use",
			sportID:       1,
			eventCount:    5,
			countError:    nil,
			deleteError:   nil,
			expectedError: true,
		},
		{
			name:          "count error",
			sportID:       1,
			eventCount:    0,
			countError:    errors.New("database error"),
			deleteError:   nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockSportRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForSport)

			service := NewSportService(mockRepo, mockEventRepo)

			mockEventRepo.On("CountEventsBySportID", mock.Anything, tt.sportID).Return(tt.eventCount, tt.countError)

			if tt.countError == nil && tt.eventCount == 0 {
				mockRepo.On("DeleteSport", mock.Anything, tt.sportID).Return(tt.deleteError)
			}

			err := service.DeleteSport(context.Background(), tt.sportID)

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

