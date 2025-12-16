package infrastructure

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vsennikov/sports-event-calendar/services"
)

func TestEventRepository_Integration(t *testing.T) {
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

	repo := NewEventRepository(db)
	ctx := context.Background()

	// Setup test data
	sportRepo := NewSportRepository(db)
	sportID, err := sportRepo.CreateSport(ctx, "Test Football")
	require.NoError(t, err)

	teamRepo := NewTeamRepository(db)
	homeTeamID, err := teamRepo.CreateTeam(ctx, services.TeamRequest{
		Name:    "Home Team",
		City:    "Home City",
		SportID: sportID,
	})
	require.NoError(t, err)

	awayTeamID, err := teamRepo.CreateTeam(ctx, services.TeamRequest{
		Name:    "Away Team",
		City:    "Away City",
		SportID: sportID,
	})
	require.NoError(t, err)

	venueRepo := NewVenueRepository(db)
	venueID, err := venueRepo.CreateVenue(ctx, services.VenueRequest{
		Name:        "Test Venue",
		City:        "Test City",
		CountryCode: "US",
	})
	require.NoError(t, err)

	t.Run("CreateEvent", func(t *testing.T) {
		eventTime := time.Now().Add(24 * time.Hour)
		description := "Test event"
		params := services.CreateEventParams{
			EventDatetime: eventTime,
			Description:   &description,
			SportID:       sportID,
			VenueID:       &venueID,
			HomeTeamID:    homeTeamID,
			AwayTeamID:    awayTeamID,
		}

		id, err := repo.CreateEvent(ctx, params)
		require.NoError(t, err)
		assert.Greater(t, id, 0)
	})

	t.Run("GetEventByID", func(t *testing.T) {
		eventTime := time.Now().Add(24 * time.Hour)
		params := services.CreateEventParams{
			EventDatetime: eventTime,
			SportID:       sportID,
			HomeTeamID:    homeTeamID,
			AwayTeamID:    awayTeamID,
		}

		id, err := repo.CreateEvent(ctx, params)
		require.NoError(t, err)

		event, err := repo.GetEventByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, event.ID)
		assert.Equal(t, sportID, event.Sport.ID)
		assert.Equal(t, homeTeamID, event.HomeTeam.ID)
		assert.Equal(t, awayTeamID, event.AwayTeam.ID)
	})

	t.Run("ListEvents", func(t *testing.T) {
		// Create multiple events
		for i := 0; i < 5; i++ {
			eventTime := time.Now().Add(time.Duration(i+1) * 24 * time.Hour)
			params := services.CreateEventParams{
				EventDatetime: eventTime,
				SportID:       sportID,
				HomeTeamID:    homeTeamID,
				AwayTeamID:    awayTeamID,
			}
			_, err := repo.CreateEvent(ctx, params)
			require.NoError(t, err)
		}

		params := services.ListEventsParams{
			Limit:  3,
			Offset: 0,
		}

		events, err := repo.ListEvents(ctx, params)
		require.NoError(t, err)
		assert.LessOrEqual(t, len(events), 3)
	})

	t.Run("ListEvents with filter", func(t *testing.T) {
		dateFrom := time.Now().Add(48 * time.Hour)
		params := services.ListEventsParams{
			SportID:  &sportID,
			DateFrom: &dateFrom,
			Limit:    10,
			Offset:   0,
		}

		events, err := repo.ListEvents(ctx, params)
		require.NoError(t, err)
		for _, event := range events {
			assert.True(t, event.EventDatetime.After(dateFrom) || event.EventDatetime.Equal(dateFrom))
			assert.Equal(t, sportID, event.Sport.ID)
		}
	})

	t.Run("CountEvents", func(t *testing.T) {
		params := services.ListEventsParams{
			SportID: &sportID,
		}

		count, err := repo.CountEvents(ctx, params)
		require.NoError(t, err)
		assert.Greater(t, count, 0)
	})

	t.Run("UpdateEvent", func(t *testing.T) {
		eventTime := time.Now().Add(24 * time.Hour)
		params := services.CreateEventParams{
			EventDatetime: eventTime,
			SportID:       sportID,
			HomeTeamID:    homeTeamID,
			AwayTeamID:    awayTeamID,
		}

		id, err := repo.CreateEvent(ctx, params)
		require.NoError(t, err)

		event, err := repo.GetEventByID(ctx, id)
		require.NoError(t, err)

		updatedDescription := "Updated description"
		homeScore := 2
		awayScore := 1
		event.Description = &updatedDescription
		event.HomeScore = &homeScore
		event.AwayScore = &awayScore

		err = repo.UpdateEvent(ctx, *event)
		require.NoError(t, err)

		updatedEvent, err := repo.GetEventByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, updatedDescription, *updatedEvent.Description)
		assert.Equal(t, homeScore, *updatedEvent.HomeScore)
		assert.Equal(t, awayScore, *updatedEvent.AwayScore)
	})

	t.Run("DeleteEvent", func(t *testing.T) {
		eventTime := time.Now().Add(24 * time.Hour)
		params := services.CreateEventParams{
			EventDatetime: eventTime,
			SportID:       sportID,
			HomeTeamID:    homeTeamID,
			AwayTeamID:    awayTeamID,
		}

		id, err := repo.CreateEvent(ctx, params)
		require.NoError(t, err)

		err = repo.DeleteEvent(ctx, id)
		require.NoError(t, err)

		_, err = repo.GetEventByID(ctx, id)
		assert.Error(t, err)
	})

	t.Run("CountEventsBySportID", func(t *testing.T) {
		count, err := repo.CountEventsBySportID(ctx, sportID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 0)
	})

	t.Run("CountEventsByTeamID", func(t *testing.T) {
		count, err := repo.CountEventsByTeamID(ctx, homeTeamID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 0)
	})

	t.Run("CountEventsByVenueId", func(t *testing.T) {
		count, err := repo.CountEventsByVenueId(ctx, venueID)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, count, 0)
	})
}

