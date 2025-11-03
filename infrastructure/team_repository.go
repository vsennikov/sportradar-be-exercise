package infrastructure

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type TeamRepository struct {
		db *sqlx.DB
}

func NewTeamRepository(db *sqlx.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(ctx context.Context, params services.TeamRequest) (int, error) {
	query := `INSERT INTO teams (name, city, _sport_id) VALUES ($1, $2, $3) RETURNING id`
	var newID int

	err := r.db.QueryRowContext(ctx, query, params.Name, params.City, params.SportID).Scan(&newID)
	if err != nil {
		return 0, err
	}
	return newID, nil
}

func (r *TeamRepository) GetTeamByID(ctx context.Context, id int) (*services.Team, error) {
	query := `SELECT id, name, city, _sport_id FROM teams WHERE id = $1`
	var dbTeam teamDBModel

	if err := r.db.GetContext(ctx, &dbTeam, query, id); err != nil {
		return nil, err
	}
	team := toServiceTeam(dbTeam)
	return &team, nil
}

func (r *TeamRepository) ListTeams(ctx context.Context) ([]services.Team, error) {
	query := `SELECT id, name, city, _sport_id FROM teams ORDER BY name ASC`
	var dbTeams []teamDBModel

	if err := r.db.SelectContext(ctx, &dbTeams, query); err != nil {
		return nil, err
	}
	teams := make([]services.Team, 0, len(dbTeams))
	for _, dbTeam := range dbTeams {
		teams = append(teams, toServiceTeam(dbTeam))
	}
	return teams, nil
}

func (r *TeamRepository) UpdateTeam(ctx context.Context, team services.Team) error {
	query := `UPDATE teams SET name = $1, city = $2 WHERE id = $3`

	_, err := r.db.ExecContext(ctx, query, team.Name, team.City, team.ID)
	return err
}

func (r *TeamRepository) DeleteTeam(ctx context.Context, id int) error {
	query := `DELETE FROM teams WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	return err
}
