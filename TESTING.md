# Testing Documentation

This document provides comprehensive information about the test suite for the Sportradar Backend Exercise project.

## Table of Contents

- [Overview](#overview)
- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Test Setup](#test-setup)
- [Unit Tests](#unit-tests)
- [Integration Tests](#integration-tests)
- [Writing New Tests](#writing-new-tests)
- [Best Practices](#best-practices)
- [Troubleshooting](#troubleshooting)

## Overview

The test suite is organized into three main categories:

1. **Unit Tests** - Fast, isolated tests using mocks
2. **Integration Tests** - Tests that interact with a real database
3. **Handler Tests** - HTTP layer tests using mocks

### Test Coverage

- ✅ Service layer business logic
- ✅ Validation rules
- ✅ Error handling
- ✅ HTTP request/response handling
- ✅ Database operations
- ✅ Foreign key relationships

## Test Structure

```
.
├── services/
│   ├── event_service_test.go      # EventService unit tests
│   ├── sport_service_test.go      # SportService unit tests
│   ├── team_service_test.go       # TeamService unit tests
│   └── venue_service_test.go      # VenueService unit tests
├── controllers/
│   ├── event_handler_test.go      # EventHandler HTTP tests
│   └── sport_handler_test.go      # SportHandler HTTP tests
└── infrastructure/
    ├── test_helpers.go                    # Test utilities
    ├── event_db_integration_test.go       # EventRepository integration tests
    ├── sport_db_integration_test.go       # SportRepository integration tests
    ├── team_repository_integration_test.go # TeamRepository integration tests
    └── venue_db_integration_test.go       # VenueRepository integration tests
```

## Running Tests

### Run All Tests

```bash
go test ./...
```

### Run Only Unit Tests (Skip Integration Tests)

```bash
go test -short ./...
```

### Run Tests for a Specific Package

```bash
# Run service tests only
go test ./services/...

# Run controller tests only
go test ./controllers/...

# Run integration tests only
go test ./infrastructure/...
```

### Run a Specific Test

```bash
# Run a specific test function
go test -run TestEventService_GetEventByID ./services/...

# Run tests matching a pattern
go test -run TestEventService ./services/...
```

### Run Tests with Verbose Output

```bash
go test -v ./...
```

### Run Tests with Coverage

```bash
# Generate coverage report
go test -cover ./...

# Generate detailed coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Run Tests in Parallel

```bash
go test -parallel 4 ./...
```

## Test Setup

### Unit Tests

Unit tests require **no setup** - they use mocks and run independently.

### Integration Tests

Integration tests require a PostgreSQL database. The tests will automatically skip if the database is not available.

#### Option 1: Using Docker Compose (Recommended)

If you have a `docker-compose.yml` with a PostgreSQL service:

```bash
# Start the database
docker compose up -d postgres

# Run integration tests
go test ./infrastructure/...
```

#### Option 2: Using Environment Variables

Set the following environment variables to configure the test database:

```bash
export TEST_DB_HOST=localhost
export TEST_DB_PORT=5432
export TEST_DB_USER=postgres
export TEST_DB_PASSWORD=postgres
export TEST_DB_NAME=sportradar_test

# Run integration tests
go test ./infrastructure/...
```

#### Option 3: Using Default Values

If no environment variables are set, the tests use these defaults:
- Host: `localhost`
- Port: `5432`
- User: `postgres`
- Password: `postgres`
- Database: `sportradar_test`

**Note:** The test database must exist. Create it manually:

```sql
CREATE DATABASE sportradar_test;
```

## Unit Tests

### Service Layer Tests

Service layer tests use mocks to isolate business logic from database operations.

#### EventService Tests (`services/event_service_test.go`)

Tests cover:
- ✅ Getting events by ID (success and not found)
- ✅ Creating events (with validation for past dates)
- ✅ Listing events (pagination, filtering, empty results)
- ✅ Updating events (partial updates, foreign key validation)
- ✅ Deleting events (success and not found)

**Key Test Cases:**
- `TestEventService_GetEventByID` - Retrieves event or handles not found
- `TestEventService_CreateEvent` - Validates past date prevention
- `TestEventService_ListEvents` - Tests pagination defaults and filtering
- `TestEventService_UpdateEvent` - Validates foreign key relationships
- `TestEventService_DeleteEvent` - Handles deletion and errors

#### SportService Tests (`services/sport_service_test.go`)

Tests cover:
- ✅ Creating sports (name length validation)
- ✅ Getting sports by ID
- ✅ Listing all sports
- ✅ Updating sports
- ✅ Deleting sports (prevents deletion when in use)

**Key Test Cases:**
- `TestSportService_CreateSport` - Validates name must be at least 3 characters
- `TestSportService_DeleteSport` - Prevents deletion when sport has events

#### TeamService Tests (`services/team_service_test.go`)

Tests cover:
- ✅ Creating teams (name and city validation)
- ✅ Getting teams by ID
- ✅ Listing all teams
- ✅ Updating teams
- ✅ Deleting teams (prevents deletion when in use)

**Key Test Cases:**
- `TestTeamService_CreateTeam` - Validates name and city must be at least 3 characters
- `TestTeamService_DeleteTeam` - Prevents deletion when team has events

#### VenueService Tests (`services/venue_service_test.go`)

Tests cover:
- ✅ Creating venues (name, city, and country code validation)
- ✅ Getting venues by ID
- ✅ Listing all venues
- ✅ Updating venues
- ✅ Deleting venues (prevents deletion when in use)

**Key Test Cases:**
- `TestVenueService_CreateVenue` - Validates name/city (3+ chars) and country code (exactly 2 chars)
- `TestVenueService_UpdateVenue` - Validates all fields during update
- `TestVenueService_DeleteVenue` - Prevents deletion when venue has events

### Handler/Controller Tests

Handler tests verify HTTP request/response handling using mocked services.

#### EventHandler Tests (`controllers/event_handler_test.go`)

Tests cover:
- ✅ POST `/events` - Creating events via HTTP
- ✅ GET `/events/:id` - Getting events by ID
- ✅ GET `/events` - Listing events with query parameters
- ✅ PATCH `/events/:id` - Updating events
- ✅ DELETE `/events/:id` - Deleting events

**Key Test Cases:**
- `TestEventHandler_HandleCreateEvent` - Validates JSON parsing and error responses
- `TestEventHandler_HandleGetEventByID` - Validates ID format and not found handling
- `TestEventHandler_HandleListEvents` - Tests query parameter parsing
- `TestEventHandler_HandleUpdateEvent` - Validates partial updates
- `TestEventHandler_HandleDeleteEvent` - Tests deletion via HTTP

#### SportHandler Tests (`controllers/sport_handler_test.go`)

Tests cover:
- ✅ POST `/sports` - Creating sports
- ✅ GET `/sports/:id` - Getting sports by ID
- ✅ GET `/sports` - Listing all sports
- ✅ PUT `/sports/:id` - Updating sports
- ✅ DELETE `/sports/:id` - Deleting sports

**Key Test Cases:**
- `TestSportHandler_HandleCreateSport` - Validates request body parsing
- `TestSportHandler_HandleDeleteSport` - Tests error handling when sport is in use

## Integration Tests

Integration tests verify database operations using a real PostgreSQL database.

### Test Helpers (`infrastructure/test_helpers.go`)

The test helpers provide:
- `SetupTestDB()` - Creates database connection (skips if unavailable)
- `InitTestSchema()` - Creates test database schema
- `CleanupTestDB()` - Cleans up test data and resets sequences

### EventRepository Integration Tests (`infrastructure/event_db_integration_test.go`)

Tests verify:
- ✅ Creating events with all relationships
- ✅ Retrieving events with joined data (sport, teams, venue)
- ✅ Listing events with pagination
- ✅ Filtering events by sport and date
- ✅ Counting events with filters
- ✅ Updating events (scores, description)
- ✅ Deleting events
- ✅ Counting events by relationships

**Key Test Cases:**
- `TestEventRepository_Integration/CreateEvent` - Tests INSERT with foreign keys
- `TestEventRepository_Integration/GetEventByID` - Tests JOIN queries
- `TestEventRepository_Integration/ListEvents` - Tests pagination at DB level
- `TestEventRepository_Integration/ListEvents with filter` - Tests WHERE clauses

### SportRepository Integration Tests (`infrastructure/sport_db_integration_test.go`)

Tests verify:
- ✅ Creating sports
- ✅ Retrieving sports
- ✅ Listing sports with ordering
- ✅ Updating sports
- ✅ Deleting sports

### TeamRepository Integration Tests (`infrastructure/team_repository_integration_test.go`)

Tests verify:
- ✅ Creating teams with sport relationship
- ✅ Retrieving teams
- ✅ Listing teams with ordering
- ✅ Updating teams
- ✅ Deleting teams

### VenueRepository Integration Tests (`infrastructure/venue_db_integration_test.go`)

Tests verify:
- ✅ Creating venues
- ✅ Retrieving venues
- ✅ Listing venues with ordering
- ✅ Updating venues
- ✅ Deleting venues

## Writing New Tests

### Adding a New Service Test

1. Create a mock for the repository interface
2. Write test cases covering success and error scenarios
3. Use table-driven tests for multiple scenarios

**Example:**

```go
func TestMyService_MyMethod(t *testing.T) {
    tests := []struct {
        name          string
        input         MyInput
        mockReturn    MyOutput
        mockError     error
        expectedError bool
    }{
        {
            name:          "successful case",
            input:         MyInput{Value: "test"},
            mockReturn:    MyOutput{Result: "success"},
            mockError:      nil,
            expectedError: false,
        },
        {
            name:          "error case",
            input:         MyInput{Value: "invalid"},
            mockReturn:    MyOutput{},
            mockError:      errors.New("validation error"),
            expectedError: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(MockMyRepository)
            service := NewMyService(mockRepo)

            mockRepo.On("MyMethod", mock.Anything, tt.input).Return(tt.mockReturn, tt.mockError)

            result, err := service.MyMethod(context.Background(), tt.input)

            if tt.expectedError {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.mockReturn, result)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

### Adding a New Handler Test

1. Create a mock for the service interface
2. Test HTTP request/response handling
3. Verify status codes and response formats

**Example:**

```go
func TestMyHandler_HandleMyEndpoint(t *testing.T) {
    tests := []struct {
        name           string
        requestBody    interface{}
        mockReturn     MyOutput
        mockError      error
        expectedStatus int
    }{
        {
            name:           "successful request",
            requestBody:    MyRequest{Value: "test"},
            mockReturn:     MyOutput{Result: "success"},
            mockError:      nil,
            expectedStatus: http.StatusOK,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockService := new(MockMyService)
            handler := NewMyHandler(mockService)

            router := setupRouter()
            router.POST("/my-endpoint", handler.HandleMyEndpoint)

            body, _ := json.Marshal(tt.requestBody)
            req := httptest.NewRequest("POST", "/my-endpoint", bytes.NewBuffer(body))
            req.Header.Set("Content-Type", "application/json")
            w := httptest.NewRecorder()

            mockService.On("MyMethod", mock.Anything, mock.Anything).Return(tt.mockReturn, tt.mockError)

            router.ServeHTTP(w, req)

            assert.Equal(t, tt.expectedStatus, w.Code)
            mockService.AssertExpectations(t)
        })
    }
}
```

### Adding a New Integration Test

1. Use `SetupTestDB()` to get a database connection
2. Use `InitTestSchema()` to set up the schema
3. Use `CleanupTestDB()` to clean up after tests
4. Test actual database operations

**Example:**

```go
func TestMyRepository_Integration(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test")
    }

    db := SetupTestDB(t)
    if db == nil {
        return
    }
    defer db.Close()

    InitTestSchema(t, db)
    defer CleanupTestDB(t, db)

    repo := NewMyRepository(db)
    ctx := context.Background()

    t.Run("CreateMyEntity", func(t *testing.T) {
        id, err := repo.CreateMyEntity(ctx, MyEntity{Name: "Test"})
        require.NoError(t, err)
        assert.Greater(t, id, 0)
    })
}
```

## Best Practices

### 1. Use Table-Driven Tests

Table-driven tests make it easy to add new test cases and keep tests organized:

```go
tests := []struct {
    name          string
    input         Input
    expected      Output
    expectedError bool
}{
    // test cases
}
```

### 2. Use Descriptive Test Names

Test names should clearly describe what is being tested:

```go
// Good
TestEventService_CreateEvent_RejectsPastDate

// Bad
TestCreateEvent1
```

### 3. Test Both Success and Error Cases

Always test:
- ✅ Happy path (successful operations)
- ✅ Validation failures
- ✅ Not found scenarios
- ✅ Database errors
- ✅ Edge cases

### 4. Use Mocks for Unit Tests

Unit tests should be fast and isolated:
- Use mocks for dependencies
- Don't make real database calls
- Don't make real HTTP requests

### 5. Use Real Database for Integration Tests

Integration tests should verify:
- Actual SQL queries work correctly
- Foreign key constraints are enforced
- Database relationships are correct

### 6. Clean Up Test Data

Always clean up test data:
- Use `defer CleanupTestDB()` in integration tests
- Reset sequences after cleanup
- Don't leave test data in the database

### 7. Skip Integration Tests When Appropriate

Use `testing.Short()` to skip integration tests:

```go
if testing.Short() {
    t.Skip("Skipping integration test")
}
```

### 8. Use Assertions from testify

Use `assert` and `require` from testify:
- `assert` - continues test execution on failure
- `require` - stops test execution on failure (use for setup)

## Troubleshooting

### Integration Tests Fail with "role does not exist"

**Problem:** PostgreSQL user doesn't exist.

**Solution:**
1. Create the user: `CREATE USER postgres WITH PASSWORD 'postgres';`
2. Or set `TEST_DB_USER` to an existing user
3. Or use Docker Compose to set up the database

### Integration Tests Skip Automatically

**Problem:** Tests skip because database connection fails.

**Solution:**
1. Ensure PostgreSQL is running
2. Check environment variables are set correctly
3. Verify database exists: `CREATE DATABASE sportradar_test;`
4. Check connection credentials

### Tests Fail with "mock expectations not met"

**Problem:** Mock setup doesn't match actual calls.

**Solution:**
1. Check mock setup matches service calls exactly
2. Verify all expected calls are made
3. Use `mock.Anything` for flexible matching

### Tests Fail with "invalid memory address"

**Problem:** Nil pointer dereference in code.

**Solution:**
1. Check for nil checks before dereferencing
2. Verify all required fields are set in test data
3. Review error handling in service code

### Coverage Report Shows Low Coverage

**Problem:** Some code paths aren't tested.

**Solution:**
1. Review coverage report: `go tool cover -html=coverage.out`
2. Add tests for untested code paths
3. Focus on error handling and edge cases

## Test Dependencies

The test suite uses the following dependencies:

- `github.com/stretchr/testify/assert` - Assertions
- `github.com/stretchr/testify/require` - Required assertions
- `github.com/stretchr/testify/mock` - Mocking framework

These are already included in `go.mod`.

## Continuous Integration

To run tests in CI/CD:

```yaml
# Example GitHub Actions workflow
- name: Run unit tests
  run: go test -short ./...

- name: Run integration tests
  run: go test ./infrastructure/...
  env:
    TEST_DB_HOST: localhost
    TEST_DB_PORT: 5432
    TEST_DB_USER: postgres
    TEST_DB_PASSWORD: postgres
    TEST_DB_NAME: sportradar_test
```

## Additional Resources

- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [testify Documentation](https://github.com/stretchr/testify)
- [Table-Driven Tests](https://github.com/golang/go/wiki/TableDrivenTests)

