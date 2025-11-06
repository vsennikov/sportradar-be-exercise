package infrastructure

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSportRepository_Integration(t *testing.T) {
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

	repo := NewSportRepository(db)
	ctx := context.Background()

	t.Run("CreateSport", func(t *testing.T) {
		id, err := repo.CreateSport(ctx, "Basketball")
		require.NoError(t, err)
		assert.Greater(t, id, 0)
	})

	t.Run("GetSportById", func(t *testing.T) {
		id, err := repo.CreateSport(ctx, "Tennis")
		require.NoError(t, err)

		sport, err := repo.GetSportById(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, id, sport.ID)
		assert.Equal(t, "Tennis", sport.Name)
	})

	t.Run("ListSports", func(t *testing.T) {
		_, err := repo.CreateSport(ctx, "Soccer")
		require.NoError(t, err)

		sports, err := repo.ListSports(ctx)
		require.NoError(t, err)
		assert.Greater(t, len(sports), 0)
	})

	t.Run("UpdateSport", func(t *testing.T) {
		id, err := repo.CreateSport(ctx, "Baseball")
		require.NoError(t, err)

		err = repo.UpdateSport(ctx, id, "Updated Baseball")
		require.NoError(t, err)

		sport, err := repo.GetSportById(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, "Updated Baseball", sport.Name)
	})

	t.Run("DeleteSport", func(t *testing.T) {
		id, err := repo.CreateSport(ctx, "Volleyball")
		require.NoError(t, err)

		err = repo.DeleteSport(ctx, id)
		require.NoError(t, err)

		_, err = repo.GetSportById(ctx, id)
		assert.Error(t, err)
	})
}

