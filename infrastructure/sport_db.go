package infrastructure

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type SportRepository struct {
	db *sqlx.DB
}

func NewSportRepository(db *sqlx.DB) *SportRepository {
	return &SportRepository{
		db: db,
	}
}

func (r *SportRepository) CreateSport(ctx context.Context, name string) (int, error) {
	query := "INSERT INTO sports (name) VALUES ($1) RETURNING id"
	var newID int

	if err := r.db.QueryRowContext(ctx, query, name).Scan(&newID); err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *SportRepository) GetSportById(ctx context.Context, id int) (*services.Sport, error) {
	query := "SELECT id, name FROM sports WHERE id = $1"
	var dbModel sportDBModel

	if err := r.db.GetContext(ctx, &dbModel, query, id); err != nil {
		return nil, err
	}
	sport := toServiceSport(dbModel)
	return &sport, nil
}

func (r *SportRepository) ListSports(ctx context.Context) ([]services.Sport, error) {
	query := "SELECT id, name FROM sports ORDER BY name ASC"
	var dbModel []sportDBModel
	if err := r.db.SelectContext(ctx, &dbModel, query); err != nil {
		return nil, err
	}
	sports := make([]services.Sport, 0, len(dbModel))
	for _, dbSport := range dbModel {
		sports = append(sports, toServiceSport(dbSport))
	}
	return sports, nil
}

func (r *SportRepository) UpdateSport(ctx context.Context, id int, name string) error {
	query := "UPDATE sports SET name = $1 WHERE id = $2"

	_, err := r.db.ExecContext(ctx, query, name, id)
	return err
}

func (r *SportRepository) DeleteSport(ctx context.Context, id int) error {
	query := "DELETE FROM sports WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
