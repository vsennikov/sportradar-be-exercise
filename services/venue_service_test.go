package services

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockVenueRepositoryForService is a mock for VenueRepositoryInterface
type MockVenueRepositoryForService struct {
	mock.Mock
}

func (m *MockVenueRepositoryForService) CreateVenue(ctx context.Context, params VenueRequest) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockVenueRepositoryForService) GetVenueById(ctx context.Context, id int) (*Venue, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Venue), args.Error(1)
}

func (m *MockVenueRepositoryForService) ListVenues(ctx context.Context) ([]Venue, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Venue), args.Error(1)
}

func (m *MockVenueRepositoryForService) UpdateVenue(ctx context.Context, venue Venue) error {
	args := m.Called(ctx, venue)
	return args.Error(0)
}

func (m *MockVenueRepositoryForService) DeleteVenue(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockEventRepositoryForVenue is a mock for EventRepositoryInterface used in VenueService tests
type MockEventRepositoryForVenue struct {
	mock.Mock
}

func (m *MockEventRepositoryForVenue) GetEventByID(ctx context.Context, id int) (*Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Event), args.Error(1)
}

func (m *MockEventRepositoryForVenue) CreateEvent(ctx context.Context, params CreateEventParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForVenue) CountEvents(ctx context.Context, params ListEventsParams) (int, error) {
	args := m.Called(ctx, params)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForVenue) ListEvents(ctx context.Context, params ListEventsParams) ([]Event, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Event), args.Error(1)
}

func (m *MockEventRepositoryForVenue) CountEventsBySportID(ctx context.Context, sportID int) (int, error) {
	args := m.Called(ctx, sportID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForVenue) CountEventsByVenueId(ctx context.Context, venueID int) (int, error) {
	args := m.Called(ctx, venueID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForVenue) CountEventsByTeamID(ctx context.Context, teamID int) (int, error) {
	args := m.Called(ctx, teamID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepositoryForVenue) UpdateEvent(ctx context.Context, event Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepositoryForVenue) DeleteEvent(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestVenueService_CreateVenue(t *testing.T) {
	tests := []struct {
		name          string
		request       CreateVenueRequest
		mockID        int
		mockError     error
		expectedID    int
		expectedError bool
	}{
		{
			name:          "successful creation",
			request:      CreateVenueRequest{Name: "Staples Center", City: "Los Angeles", CountryCode: "US"},
			mockID:        1,
			mockError:     nil,
			expectedID:    1,
			expectedError: false,
		},
		{
			name:          "name too short",
			request:      CreateVenueRequest{Name: "AB", City: "Los Angeles", CountryCode: "US"},
			mockID:        0,
			mockError:     nil,
			expectedID:    0,
			expectedError: true,
		},
		{
			name:          "city too short",
			request:      CreateVenueRequest{Name: "Staples Center", City: "LA", CountryCode: "US"},
			mockID:        0,
			mockError:     nil,
			expectedID:    0,
			expectedError: true,
		},
		{
			name:          "country code wrong length",
			request:      CreateVenueRequest{Name: "Staples Center", City: "Los Angeles", CountryCode: "USA"},
			mockID:        0,
			mockError:     nil,
			expectedID:    0,
			expectedError: true,
		},
		{
			name:          "database error",
			request:      CreateVenueRequest{Name: "Staples Center", City: "Los Angeles", CountryCode: "US"},
			mockID:        0,
			mockError:     errors.New("database error"),
			expectedID:    0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVenueRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForVenue)

			service := NewVenueService(mockRepo, mockEventRepo)

			if !tt.expectedError || tt.name == "database error" {
				mockRepo.On("CreateVenue", mock.Anything, mock.AnythingOfType("VenueRequest")).Return(tt.mockID, tt.mockError)
			}

			result, err := service.CreateVenue(context.Background(), tt.request)

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

func TestVenueService_GetVenueByID(t *testing.T) {
	tests := []struct {
		name          string
		venueID       int
		mockVenue     *Venue
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful retrieval",
			venueID:       1,
			mockVenue:     &Venue{ID: 1, Name: "Staples Center", City: "Los Angeles", CountryCode: "US"},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "venue not found",
			venueID:       999,
			mockVenue:     nil,
			mockError:     sql.ErrNoRows,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVenueRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForVenue)

			service := NewVenueService(mockRepo, mockEventRepo)

			mockRepo.On("GetVenueById", mock.Anything, tt.venueID).Return(tt.mockVenue, tt.mockError)

			result, err := service.GetVenueByID(context.Background(), tt.venueID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.venueID, result.ID)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestVenueService_ListVenues(t *testing.T) {
	tests := []struct {
		name          string
		mockVenues    []Venue
		mockError     error
		expectedCount int
		expectedError bool
	}{
		{
			name:          "successful list",
			mockVenues:    []Venue{{ID: 1, Name: "Staples Center"}, {ID: 2, Name: "Madison Square Garden"}},
			mockError:     nil,
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:          "empty list",
			mockVenues:    []Venue{},
			mockError:     nil,
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:          "database error",
			mockVenues:    nil,
			mockError:     errors.New("database error"),
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVenueRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForVenue)

			service := NewVenueService(mockRepo, mockEventRepo)

			mockRepo.On("ListVenues", mock.Anything).Return(tt.mockVenues, tt.mockError)

			result, err := service.ListVenues(context.Background())

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

func TestVenueService_UpdateVenue(t *testing.T) {
	existingVenue := &Venue{
		ID:          1,
		Name:        "Staples Center",
		City:        "Los Angeles",
		CountryCode: "US",
	}

	tests := []struct {
		name          string
		venueID       int
		request       UpdateVenueRequest
		mockVenue     *Venue
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful update",
			venueID:       1,
			request:       UpdateVenueRequest{Name: stringPtr("Updated Venue")},
			mockVenue:     existingVenue,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "name too short",
			venueID:       1,
			request:       UpdateVenueRequest{Name: stringPtr("AB")},
			mockVenue:     existingVenue,
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "city too short",
			venueID:       1,
			request:       UpdateVenueRequest{City: stringPtr("LA")},
			mockVenue:     existingVenue,
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "country code wrong length",
			venueID:       1,
			request:       UpdateVenueRequest{CountryCode: stringPtr("USA")},
			mockVenue:     existingVenue,
			mockError:     nil,
			expectedError: true,
		},
		{
			name:          "venue not found",
			venueID:       999,
			request:       UpdateVenueRequest{Name: stringPtr("Updated")},
			mockVenue:     nil,
			mockError:     sql.ErrNoRows,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVenueRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForVenue)

			service := NewVenueService(mockRepo, mockEventRepo)

			mockRepo.On("GetVenueById", mock.Anything, tt.venueID).Return(tt.mockVenue, tt.mockError)

			if tt.mockVenue != nil && !tt.expectedError {
				mockRepo.On("UpdateVenue", mock.Anything, mock.AnythingOfType("Venue")).Return(nil)
			}

			err := service.UpdateVenue(context.Background(), tt.venueID, tt.request)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestVenueService_DeleteVenue(t *testing.T) {
	tests := []struct {
		name          string
		venueID       int
		eventCount    int
		countError    error
		deleteError   error
		expectedError bool
	}{
		{
			name:          "successful deletion",
			venueID:       1,
			eventCount:    0,
			countError:    nil,
			deleteError:   nil,
			expectedError: false,
		},
		{
			name:          "venue in use",
			venueID:       1,
			eventCount:    5,
			countError:    nil,
			deleteError:   nil,
			expectedError: true,
		},
		{
			name:          "count error",
			venueID:       1,
			eventCount:    0,
			countError:    errors.New("database error"),
			deleteError:   nil,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockVenueRepositoryForService)
			mockEventRepo := new(MockEventRepositoryForVenue)

			service := NewVenueService(mockRepo, mockEventRepo)

			mockEventRepo.On("CountEventsByVenueId", mock.Anything, tt.venueID).Return(tt.eventCount, tt.countError)

			if tt.countError == nil && tt.eventCount == 0 {
				mockRepo.On("DeleteVenue", mock.Anything, tt.venueID).Return(tt.deleteError)
			}

			err := service.DeleteVenue(context.Background(), tt.venueID)

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

