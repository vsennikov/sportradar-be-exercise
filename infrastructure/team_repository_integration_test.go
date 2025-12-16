package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vsennikov/sports-event-calendar/services"
)

func TestTeamRepository_Integration(t *testing.T) {
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

	sportRepo := NewSportRepository(db)
	sportID, err := sportRepo.CreateSport(context.Background(), "Test Sport")
	require.NoError(t, err)

	repo := NewTeamRepository(db)
	ctx := context.Background()

	t.Run("CreateTeam", func(t *testing.T) {
		params := services.TeamRequest{
			Name:    "Lakers",
			City:    "Los Angeles",
			SportID: sportID,
		}

		id, err := repo.CreateTeam(ctx, params)
		require.NoError(t, err)
		assert.Greater(t, id, 0)
	})

	t.Run("GetTeamByID", func(t *testing.T) {
		params := services.TeamRequest{
			Name:    "Warriors",
			City:    "San Francisco",
			SportID: sportID,
		}

		id, err := repo.CreateTeam(ctx, params)
		require.NoError(t, err)

		team, err := repo.GetTeamByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, team.ID)
		assert.Equal(t, "Warriors", team.Name)
		assert.Equal(t, "San Francisco", team.City)
	})

	t.Run("ListTeams", func(t *testing.T) {
		params := services.TeamRequest{
			Name:    "Celtics",
			City:    "Boston",
			SportID: sportID,
		}

		_, err := repo.CreateTeam(ctx, params)
		require.NoError(t, err)

		teams, err := repo.ListTeams(ctx)
		require.NoError(t, err)
		assert.Greater(t, len(teams), 0)
	})

	t.Run("UpdateTeam", func(t *testing.T) {
		params := services.TeamRequest{
			Name:    "Heat",
			City:    "Miami",
			SportID: sportID,
		}

		id, err := repo.CreateTeam(ctx, params)
		require.NoError(t, err)

		team, err := repo.GetTeamByID(ctx, id)
		require.NoError(t, err)

		team.Name = "Updated Heat"
		team.City = "Updated Miami"

		err = repo.UpdateTeam(ctx, *team)
		require.NoError(t, err)

		updatedTeam, err := repo.GetTeamByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Heat", updatedTeam.Name)
		assert.Equal(t, "Updated Miami", updatedTeam.City)
	})

	t.Run("DeleteTeam", func(t *testing.T) {
		params := services.TeamRequest{
			Name:    "Bulls",
			City:    "Chicago",
			SportID: sportID,
		}

		id, err := repo.CreateTeam(ctx, params)
		require.NoError(t, err)

		err = repo.DeleteTeam(ctx, id)
		require.NoError(t, err)

		_, err = repo.GetTeamByID(ctx, id)
		assert.Error(t, err)
	})
}

