package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vsennikov/sports-event-calendar/services"
)

// MockSportService is a mock implementation of SportServiceInterface
type MockSportService struct {
	mock.Mock
}

func (m *MockSportService) CreateSport(ctx context.Context, req services.SportRequest) (int, error) {
	args := m.Called(ctx, req)
	return args.Int(0), args.Error(1)
}

func (m *MockSportService) GetSportByID(ctx context.Context, id int) (*services.Sport, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.Sport), args.Error(1)
}

func (m *MockSportService) ListSports(ctx context.Context) ([]services.Sport, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]services.Sport), args.Error(1)
}

func (m *MockSportService) UpdateSport(ctx context.Context, id int, req services.SportRequest) error {
	args := m.Called(ctx, id, req)
	return args.Error(0)
}

func (m *MockSportService) DeleteSport(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestSportHandler_HandleCreateSport(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    interface{}
		mockID         int
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful creation",
			requestBody:    services.SportRequest{Name: "Basketball"},
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
			name:           "service error",
			requestBody:    services.SportRequest{Name: "Basketball"},
			mockID:         0,
			mockError:      fmt.Errorf("validation error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSportService)
			handler := NewSportHandler(mockService)

			router := setupRouter()
			router.POST("/sports", handler.HandleCreateSport)

			var body []byte
			var err error
			if tt.name == "invalid request body" {
				body = []byte("invalid json")
			} else {
				body, err = json.Marshal(tt.requestBody)
				assert.NoError(t, err)
			}

			req := httptest.NewRequest("POST", "/sports", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tt.name != "invalid request body" {
				mockService.On("CreateSport", mock.Anything, mock.AnythingOfType("SportRequest")).Return(tt.mockID, tt.mockError)
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

func TestSportHandler_HandleGetSportByID(t *testing.T) {
	tests := []struct {
		name           string
		sportID        string
		mockSport      *services.Sport
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful retrieval",
			sportID:        "1",
			mockSport:      &services.Sport{ID: 1, Name: "Football"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			sportID:        "invalid",
			mockSport:      nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "sport not found",
			sportID:        "999",
			mockSport:      nil,
			mockError:      fmt.Errorf("sport with id 999 not found"),
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSportService)
			handler := NewSportHandler(mockService)

			router := setupRouter()
			router.GET("/sports/:id", handler.HandleGetSportByID)

			req := httptest.NewRequest("GET", "/sports/"+tt.sportID, nil)
			w := httptest.NewRecorder()

			if tt.name != "invalid ID format" {
				id, _ := strconv.Atoi(tt.sportID)
				mockService.On("GetSportByID", mock.Anything, id).Return(tt.mockSport, tt.mockError)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response sportDTO
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockSport.ID, response.ID)
			}

			if tt.name != "invalid ID format" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestSportHandler_HandleListSports(t *testing.T) {
	tests := []struct {
		name           string
		mockSports     []services.Sport
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful list",
			mockSports:     []services.Sport{{ID: 1, Name: "Football"}, {ID: 2, Name: "Basketball"}},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "service error",
			mockSports:     nil,
			mockError:      fmt.Errorf("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSportService)
			handler := NewSportHandler(mockService)

			router := setupRouter()
			router.GET("/sports", handler.HandleListSports)

			req := httptest.NewRequest("GET", "/sports", nil)
			w := httptest.NewRecorder()

			mockService.On("ListSports", mock.Anything).Return(tt.mockSports, tt.mockError)

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response []sportDTO
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockSports), len(response))
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestSportHandler_HandleUpdateSport(t *testing.T) {
	tests := []struct {
		name           string
		sportID        string
		requestBody    interface{}
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful update",
			sportID:        "1",
			requestBody:    services.SportRequest{Name: "Updated Football"},
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			sportID:        "invalid",
			requestBody:    services.SportRequest{Name: "Updated"},
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "service error",
			sportID:        "1",
			requestBody:    services.SportRequest{Name: "Updated"},
			mockError:      fmt.Errorf("validation error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSportService)
			handler := NewSportHandler(mockService)

			router := setupRouter()
			router.PUT("/sports/:id", handler.HandleUpdateSport)

			body, err := json.Marshal(tt.requestBody)
			assert.NoError(t, err)

			req := httptest.NewRequest("PUT", "/sports/"+tt.sportID, bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			if tt.name != "invalid ID format" {
				id, _ := strconv.Atoi(tt.sportID)
				mockService.On("UpdateSport", mock.Anything, id, mock.AnythingOfType("SportRequest")).Return(tt.mockError)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.name != "invalid ID format" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestSportHandler_HandleDeleteSport(t *testing.T) {
	tests := []struct {
		name           string
		sportID        string
		mockError      error
		expectedStatus int
	}{
		{
			name:           "successful deletion",
			sportID:        "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid ID format",
			sportID:        "invalid",
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "sport in use",
			sportID:        "1",
			mockError:      fmt.Errorf("cannot delete sport: it is currently used by 5 events"),
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := new(MockSportService)
			handler := NewSportHandler(mockService)

			router := setupRouter()
			router.DELETE("/sports/:id", handler.HandleDeleteSport)

			req := httptest.NewRequest("DELETE", "/sports/"+tt.sportID, nil)
			w := httptest.NewRecorder()

			if tt.name != "invalid ID format" {
				id, _ := strconv.Atoi(tt.sportID)
				mockService.On("DeleteSport", mock.Anything, id).Return(tt.mockError)
			}

			router.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.name != "invalid ID format" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

