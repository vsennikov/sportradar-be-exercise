# Testing Documentation

This document provides information about the test suite for the Sports Event Calendar project.

## Table of Contents

- [Overview](#overview)
- [Test Structure](#test-structure)
- [Running Tests](#running-tests)
- [Test Setup](#test-setup)
- [Unit Tests](#unit-tests)
- [Integration Tests](#integration-tests)

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
go test ./services/...

go test ./controllers/...

go test ./infrastructure/...
```

## Test Setup

### Integration Tests

Integration tests require a PostgreSQL database. The tests will automatically skip if the database is not available.

```bash
# Start the database
docker compose up -d db

# Run integration tests
go test ./infrastructure/...
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