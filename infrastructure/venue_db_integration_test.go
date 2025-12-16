package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vsennikov/sports-event-calendar/services"
)

func TestVenueRepository_Integration(t *testing.T) {
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

	repo := NewVenueRepository(db)
	ctx := context.Background()

	t.Run("CreateVenue", func(t *testing.T) {
		params := services.VenueRequest{
			Name:        "Staples Center",
			City:        "Los Angeles",
			CountryCode: "US",
		}

		id, err := repo.CreateVenue(ctx, params)
		require.NoError(t, err)
		assert.Greater(t, id, 0)
	})

	t.Run("GetVenueById", func(t *testing.T) {
		params := services.VenueRequest{
			Name:        "Madison Square Garden",
			City:        "New York",
			CountryCode: "US",
		}

		id, err := repo.CreateVenue(ctx, params)
		require.NoError(t, err)

		venue, err := repo.GetVenueById(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, venue.ID)
		assert.Equal(t, "Madison Square Garden", venue.Name)
		assert.Equal(t, "New York", venue.City)
		assert.Equal(t, "US", venue.CountryCode)
	})

	t.Run("ListVenues", func(t *testing.T) {
		params := services.VenueRequest{
			Name:        "TD Garden",
			City:        "Boston",
			CountryCode: "US",
		}

		_, err := repo.CreateVenue(ctx, params)
		require.NoError(t, err)

		venues, err := repo.ListVenues(ctx)
		require.NoError(t, err)
		assert.Greater(t, len(venues), 0)
	})

	t.Run("UpdateVenue", func(t *testing.T) {
		params := services.VenueRequest{
			Name:        "Oracle Arena",
			City:        "Oakland",
			CountryCode: "US",
		}

		id, err := repo.CreateVenue(ctx, params)
		require.NoError(t, err)

		venue, err := repo.GetVenueById(ctx, id)
		require.NoError(t, err)

		venue.Name = "Updated Oracle Arena"
		venue.City = "Updated Oakland"
		venue.CountryCode = "CA"

		err = repo.UpdateVenue(ctx, *venue)
		require.NoError(t, err)

		updatedVenue, err := repo.GetVenueById(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Oracle Arena", updatedVenue.Name)
		assert.Equal(t, "Updated Oakland", updatedVenue.City)
		assert.Equal(t, "CA", updatedVenue.CountryCode)
	})

	t.Run("DeleteVenue", func(t *testing.T) {
		params := services.VenueRequest{
			Name:        "United Center",
			City:        "Chicago",
			CountryCode: "US",
		}

		id, err := repo.CreateVenue(ctx, params)
		require.NoError(t, err)

		err = repo.DeleteVenue(ctx, id)
		require.NoError(t, err)

		_, err = repo.GetVenueById(ctx, id)
		assert.Error(t, err)
	})
}

