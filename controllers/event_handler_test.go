package controllers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

// MockEventService is a mock implementation of EventServiceInterface
type MockEventService struct {
	mock.Mock
}

func (m *MockEventService) GetEventByID(ctx context.Context, id int) (*services.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Event), args.Error(1)
}

func (m *MockEventService) CreateEvent(ctx context.Context, req services.EventCreateRequest) (int, error) {
	args := m.Called(ctx, req)
	return args.Int(0), args.Error(1)
}

func (m *MockEventService) ListEvents(ctx context.Context, req services.ListEventsRequest) ([]services.Event, *services.Pagination, error) {
	args := m.Called(ctx, req)
	var events []services.Event
	var pagination *services.Pagination
	if args.Get(0) != nil {
		events = args.Get(0).([]services.Event)
	}
	if args.Get(1) != nil {
		pagination = args.Get(1).(*services.Pagination)
	}
	return events, pagination, args.Error(2)
}

func (m *MockEventService) UpdateEvent(ctx context.Context, id int, req services.UpdateEventRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockEventService) DeleteEvent(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestEventHandler_HandleCreateEvent(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockID         int
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: services.EventCreateRequest{
				EventDatetime: time.Now().Add(24 * time.Hour),
				SportID:       1,
				HomeTeamID:    1,
				AwayTeamID:    2,
			},
			mockID:         1,
			mockError:      nil,
			expectedStatus: http.StatusCreated,
		},
		{
			name:           "invalid request body",
			requestBody:    "invalid",
			mockID:         0,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "service error",
			requestBody: services.EventCreateRequest{
				EventDatetime: time.Now().Add(24 * time.Hour),
				SportID:       1,
				HomeTeamID:    1,
				AwayTeamID:    2,
			},
			mockID:         0,
			mockError:      sql.ErrNoRows,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockEventService)
			handler := NewEventHandler(mockService)

			router := setupRouter()
			router.POST("/events", handler.HandleCreateEvent)

			var body []byte
			var err error
			if tt.name == "invalid request body" {
				body = []byte("invalid json")
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/events", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tt.name != "invalid request body" {
				mockService.On("CreateEvent", mock.Anything, mock.AnythingOfType("EventCreateRequest")).Return(tt.mockID, tt.mockError)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusCreated {
				var response map[string]int
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockID, response["id"])
			}

			if tt.name != "invalid request body" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestEventHandler_HandleGetEventByID(t *testing.T) {
	tests := []struct {
		name           string
		eventID        string
		mockEvent      *services.Event
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful retrieval",
			eventID: "1",
			mockEvent: &services.Event{
				ID:            1,
				EventDatetime: time.Now().Add(24 * time.Hour),
				Sport:         services.Sport{ID: 1, Name: "Football"},
				HomeTeam:      services.Team{ID: 1, Name: "Team A", City: "City A"},
				AwayTeam:      services.Team{ID: 2, Name: "Team B", City: "City B"},
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			eventID:        "invalid",
			mockEvent:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "event not found",
			eventID:        "999",
			mockEvent:      nil,
			mockError:      sql.ErrNoRows,
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "service error",
			eventID:        "1",
			mockEvent:      nil,
			mockError:      sql.ErrNoRows,
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockEventService)
			handler := NewEventHandler(mockService)

			router := setupRouter()
			router.GET("/events/:id", handler.HandleGetEventByID)

			req := httptest.NewRequest("GET", "/events/"+tt.eventID, nil)
			w := httptest.NewRecorder()

			if tt.name != "invalid ID format" {
				id, _ := parseID(tt.eventID)
				mockService.On("GetEventByID", mock.Anything, id).Return(tt.mockEvent, tt.mockError)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response EventDTO
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockEvent.ID, response.ID)
			}

			if tt.name != "invalid ID format" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestEventHandler_HandleListEvents(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockEvents     []services.Event
		mockPagination *services.Pagination
		mockError      error
		expectedStatus int
	}{
		{
			name:        "successful list",
			queryParams: "",
			mockEvents: []services.Event{
				{ID: 1, EventDatetime: time.Now().Add(24 * time.Hour)},
				{ID: 2, EventDatetime: time.Now().Add(48 * time.Hour)},
			},
			mockPagination: &services.Pagination{
				TotalItems:  2,
				TotalPages:  1,
				CurrentPage: 1,
				PageSize:    10,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:        "with pagination",
			queryParams: "?page=1&limit=5",
			mockEvents: []services.Event{
				{ID: 1, EventDatetime: time.Now().Add(24 * time.Hour)},
			},
			mockPagination: &services.Pagination{
				TotalItems:  1,
				TotalPages:  1,
				CurrentPage: 1,
				PageSize:    5,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "service error",
			queryParams:    "",
			mockEvents:     nil,
			mockPagination: nil,
			mockError:      sql.ErrNoRows,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockEventService)
			handler := NewEventHandler(mockService)

			router := setupRouter()
			router.GET("/events", handler.HandleListEvents)

			req := httptest.NewRequest("GET", "/events"+tt.queryParams, nil)
			w := httptest.NewRecorder()

			mockService.On("ListEvents", mock.Anything, mock.AnythingOfType("ListEventsRequest")).Return(tt.mockEvents, tt.mockPagination, tt.mockError)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.NotNil(t, response["events"])
				assert.NotNil(t, response["pagination"])
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestEventHandler_HandleUpdateEvent(t *testing.T) {
	tests := []struct {
		name           string
		eventID        string
		requestBody    interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name:    "successful update",
			eventID: "1",
			requestBody: services.UpdateEventRequest{
				Description: stringPtr("Updated description"),
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			eventID:        "invalid",
			requestBody:    services.UpdateEventRequest{},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "invalid request body",
			eventID:        "1",
			requestBody:    "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:    "service error",
			eventID: "1",
			requestBody: services.UpdateEventRequest{
				Description: stringPtr("Updated"),
			},
			mockError:      sql.ErrNoRows,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockEventService)
			handler := NewEventHandler(mockService)

			router := setupRouter()
			router.PATCH("/events/:id", handler.HandleUpdateEvent)

			var body []byte
			var err error
			if tt.name == "invalid request body" {
				body = []byte("invalid json")
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("PATCH", "/events/"+tt.eventID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tt.name != "invalid ID format" && tt.name != "invalid request body" {
				id, _ := parseID(tt.eventID)
				mockService.On("UpdateEvent", mock.Anything, id, mock.AnythingOfType("UpdateEventRequest")).Return(tt.mockError)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.name != "invalid ID format" && tt.name != "invalid request body" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestEventHandler_HandleDeleteEvent(t *testing.T) {
	tests := []struct {
		name           string
		eventID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			eventID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			eventID:        "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			eventID:        "1",
			mockError:      sql.ErrNoRows,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockEventService)
			handler := NewEventHandler(mockService)

			router := setupRouter()
			router.DELETE("/events/:id", handler.HandleDeleteEvent)

			req := httptest.NewRequest("DELETE", "/events/"+tt.eventID, nil)
			w := httptest.NewRecorder()

			if tt.name != "invalid ID format" {
				id, _ := parseID(tt.eventID)
				mockService.On("DeleteEvent", mock.Anything, id).Return(tt.mockError)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.name != "invalid ID format" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func parseID(idStr string) (int, error) {
	return strconv.Atoi(idStr)
}

