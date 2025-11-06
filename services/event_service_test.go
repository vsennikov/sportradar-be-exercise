package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockEventRepository is a mock implementation of EventRepositoryInterface
type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) GetEventByID(ctx context.Context, id int) (*Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Event), args.Error(1)
}

func (m *MockEventRepository) CreateEvent(ctx context.Context, params CreateEventParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) CountEvents(ctx context.Context, params ListEventsParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Event), args.Error(1)
}

func (m *MockEventRepository) CountEventsBySportID(ctx context.Context, sportID int) (int, error) {
	args := m.Called(ctx, sportID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) CountEventsByVenueId(ctx context.Context, venueID int) (int, error) {
	args := m.Called(ctx, venueID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) CountEventsByTeamID(ctx context.Context, teamID int) (int, error) {
	args := m.Called(ctx, teamID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) UpdateEvent(ctx context.Context, event Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) DeleteEvent(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockSportRepository is a mock implementation of SportRepositoryInterface
type MockSportRepository struct {
	mock.Mock
}

func (m *MockSportRepository) CreateSport(ctx context.Context, name string) (int, error) {
	args := m.Called(ctx, name)
	return args.Int(0), args.Error(1)
}

func (m *MockSportRepository) GetSportById(ctx context.Context, id int) (*Sport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Sport), args.Error(1)
}

func (m *MockSportRepository) ListSports(ctx context.Context) ([]Sport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Sport), args.Error(1)
}

func (m *MockSportRepository) UpdateSport(ctx context.Context, id int, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockSportRepository) DeleteSport(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockTeamRepository is a mock implementation of TeamRepositoryInterface
type MockTeamRepository struct {
	mock.Mock
}

func (m *MockTeamRepository) CreateTeam(ctx context.Context, params TeamRequest) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockTeamRepository) GetTeamByID(ctx context.Context, id int) (*Team, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Team), args.Error(1)
}

func (m *MockTeamRepository) ListTeams(ctx context.Context) ([]Team, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Team), args.Error(1)
}

func (m *MockTeamRepository) UpdateTeam(ctx context.Context, team Team) error {
	args := m.Called(ctx, team)
	return args.Error(0)
}

func (m *MockTeamRepository) DeleteTeam(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockVenueRepository is a mock implementation of VenueRepositoryInterface
type MockVenueRepository struct {
	mock.Mock
}

func (m *MockVenueRepository) CreateVenue(ctx context.Context, params VenueRequest) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockVenueRepository) GetVenueById(ctx context.Context, id int) (*Venue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Venue), args.Error(1)
}

func (m *MockVenueRepository) ListVenues(ctx context.Context) ([]Venue, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Venue), args.Error(1)
}

func (m *MockVenueRepository) UpdateVenue(ctx context.Context, venue Venue) error {
	args := m.Called(ctx, venue)
	return args.Error(0)
}

func (m *MockVenueRepository) DeleteVenue(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestEventService_GetEventByID(t *testing.T) {
	tests := []struct {
		name          string
		eventID       int
		mockEvent     *Event
		mockError     error
		expectedEvent *Event
		expectedError bool
	}{
		{
			name:    "successful retrieval",
			eventID: 1,
			mockEvent: &Event{
				ID:            1,
				EventDatetime: time.Now().Add(24 * time.Hour),
				Description:   stringPtr("Test event"),
			},
			mockError:     nil,
			expectedEvent: &Event{ID: 1},
			expectedError: false,
		},
		{
			name:          "event not found",
			eventID:       999,
			mockEvent:     nil,
			mockError:     sql.ErrNoRows,
			expectedEvent: nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			mockSportRepo := new(MockSportRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockVenueRepo := new(MockVenueRepository)

			service := NewEventService(mockRepo, 1, 10, mockSportRepo, mockTeamRepo, mockVenueRepo)

			mockRepo.On("GetEventByID", mock.Anything, tt.eventID).Return(tt.mockEvent, tt.mockError)

			result, err := service.GetEventByID(context.Background(), tt.eventID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.eventID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_CreateEvent(t *testing.T) {
	tests := []struct {
		name          string
		request       EventCreateRequest
		mockID        int
		mockError     error
		expectedID    int
		expectedError bool
	}{
		{
			name: "successful creation",
			request: EventCreateRequest{
				EventDatetime: time.Now().Add(24 * time.Hour),
				SportID:       1,
				HomeTeamID:    1,
				AwayTeamID:    2,
			},
			mockID:        1,
			mockError:     nil,
			expectedID:    1,
			expectedError: false,
		},
		{
			name: "event in the past",
			request: EventCreateRequest{
				EventDatetime: time.Now().Add(-24 * time.Hour),
				SportID:       1,
				HomeTeamID:    1,
				AwayTeamID:    2,
			},
			mockID:        0,
			mockError:     nil,
			expectedID:    0,
			expectedError: true,
		},
		{
			name: "database error",
			request: EventCreateRequest{
				EventDatetime: time.Now().Add(24 * time.Hour),
				SportID:       1,
				HomeTeamID:    1,
				AwayTeamID:    2,
			},
			mockID:        0,
			mockError:     errors.New("database error"),
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			mockSportRepo := new(MockSportRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockVenueRepo := new(MockVenueRepository)

			service := NewEventService(mockRepo, 1, 10, mockSportRepo, mockTeamRepo, mockVenueRepo)

			if !tt.expectedError || tt.name == "database error" {
				mockRepo.On("CreateEvent", mock.Anything, mock.AnythingOfType("CreateEventParams")).Return(tt.mockID, tt.mockError)
			}

			result, err := service.CreateEvent(context.Background(), tt.request)

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

func TestEventService_ListEvents(t *testing.T) {
	tests := []struct {
		name           string
		request        ListEventsRequest
		mockCount      int
		mockEvents     []Event
		mockCountError error
		mockListError  error
		expectedCount  int
		expectedError  bool
	}{
		{
			name: "successful list with pagination",
			request: ListEventsRequest{
				Page:  1,
				Limit: 10,
			},
			mockCount:      25,
			mockEvents:     []Event{{ID: 1}, {ID: 2}},
			mockCountError: nil,
			mockListError:  nil,
			expectedCount:  2,
			expectedError:  false,
		},
		{
			name: "empty result",
			request: ListEventsRequest{
				Page:  1,
				Limit: 10,
			},
			mockCount:      0,
			mockEvents:     []Event{},
			mockCountError: nil,
			mockListError:  nil,
			expectedCount:  0,
			expectedError:  false,
		},
		{
			name: "default pagination values",
			request: ListEventsRequest{
				Page:  0,
				Limit: 0,
			},
			mockCount:      5,
			mockEvents:     []Event{{ID: 1}, {ID: 2}, {ID: 3}},
			mockCountError: nil,
			mockListError:  nil,
			expectedCount:  3,
			expectedError:  false,
		},
		{
			name: "count error",
			request: ListEventsRequest{
				Page:  1,
				Limit: 10,
			},
			mockCount:      0,
			mockEvents:     nil,
			mockCountError: errors.New("database error"),
			mockListError:  nil,
			expectedCount:  0,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			mockSportRepo := new(MockSportRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockVenueRepo := new(MockVenueRepository)

			service := NewEventService(mockRepo, 1, 10, mockSportRepo, mockTeamRepo, mockVenueRepo)

			mockRepo.On("CountEvents", mock.Anything, mock.AnythingOfType("ListEventsParams")).Return(tt.mockCount, tt.mockCountError)

			if tt.mockCountError == nil && tt.mockCount > 0 {
				mockRepo.On("ListEvents", mock.Anything, mock.AnythingOfType("ListEventsParams")).Return(tt.mockEvents, tt.mockListError)
			}

			events, pagination, err := service.ListEvents(context.Background(), tt.request)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, events)
				assert.Nil(t, pagination)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, events)
				assert.NotNil(t, pagination)
				assert.Equal(t, tt.expectedCount, len(events))
				if tt.mockCount > 0 {
					assert.Equal(t, tt.mockCount, pagination.TotalItems)
					// When Page/Limit are 0, service sets defaults, so check actual pagination values
					expectedPage := tt.request.Page
					expectedLimit := tt.request.Limit
					if expectedPage <= 0 {
						expectedPage = 1 // default page
					}
					if expectedLimit <= 0 {
						expectedLimit = 10 // default limit
					}
					assert.Equal(t, expectedPage, pagination.CurrentPage)
					assert.Equal(t, expectedLimit, pagination.PageSize)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_UpdateEvent(t *testing.T) {
	existingEvent := &Event{
		ID:            1,
		EventDatetime: time.Now().Add(24 * time.Hour),
		Description:   stringPtr("Original description"),
		Sport:         Sport{ID: 1, Name: "Football"},
		HomeTeam:      Team{ID: 1, Name: "Team A"},
		AwayTeam:      Team{ID: 2, Name: "Team B"},
	}

	tests := []struct {
		name          string
		eventID       int
		request       UpdateEventRequest
		mockEvent     *Event
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful update",
			eventID:       1,
			request:       UpdateEventRequest{Description: stringPtr("Updated description")},
			mockEvent:     existingEvent,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "update sport",
			eventID:       1,
			request:       UpdateEventRequest{SportID: intPtr(2)},
			mockEvent:     existingEvent,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "event not found",
			eventID:       999,
			request:       UpdateEventRequest{Description: stringPtr("Updated")},
			mockEvent:     nil,
			mockError:     sql.ErrNoRows,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			mockSportRepo := new(MockSportRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockVenueRepo := new(MockVenueRepository)

			service := NewEventService(mockRepo, 1, 10, mockSportRepo, mockTeamRepo, mockVenueRepo)

			mockRepo.On("GetEventByID", mock.Anything, tt.eventID).Return(tt.mockEvent, tt.mockError)

			if tt.mockEvent != nil {
				if tt.request.SportID != nil {
					mockSportRepo.On("GetSportById", mock.Anything, *tt.request.SportID).Return(&Sport{ID: *tt.request.SportID}, nil)
				}
				if tt.request.VenueID != nil {
					mockVenueRepo.On("GetVenueById", mock.Anything, *tt.request.VenueID).Return(&Venue{ID: *tt.request.VenueID}, nil)
				}
				if tt.request.HomeTeamID != nil {
					mockTeamRepo.On("GetTeamByID", mock.Anything, *tt.request.HomeTeamID).Return(&Team{ID: *tt.request.HomeTeamID}, nil)
				}
				if tt.request.AwayTeamID != nil {
					mockTeamRepo.On("GetTeamByID", mock.Anything, *tt.request.AwayTeamID).Return(&Team{ID: *tt.request.AwayTeamID}, nil)
				}
				mockRepo.On("UpdateEvent", mock.Anything, mock.AnythingOfType("Event")).Return(nil)
			}

			err := service.UpdateEvent(context.Background(), tt.eventID, tt.request)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestEventService_DeleteEvent(t *testing.T) {
	tests := []struct {
		name          string
		eventID       int
		mockEvent     *Event
		mockGetError  error
		mockDelError  error
		expectedError bool
	}{
		{
			name:          "successful deletion",
			eventID:       1,
			mockEvent:     &Event{ID: 1},
			mockGetError:  nil,
			mockDelError:  nil,
			expectedError: false,
		},
		{
			name:          "event not found",
			eventID:       999,
			mockEvent:     nil,
			mockGetError:  sql.ErrNoRows,
			mockDelError:  nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockEventRepository)
			mockSportRepo := new(MockSportRepository)
			mockTeamRepo := new(MockTeamRepository)
			mockVenueRepo := new(MockVenueRepository)

			service := NewEventService(mockRepo, 1, 10, mockSportRepo, mockTeamRepo, mockVenueRepo)

			mockRepo.On("GetEventByID", mock.Anything, tt.eventID).Return(tt.mockEvent, tt.mockGetError)

			if tt.mockEvent != nil {
				mockRepo.On("DeleteEvent", mock.Anything, tt.eventID).Return(tt.mockDelError)
			}

			err := service.DeleteEvent(context.Background(), tt.eventID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}
