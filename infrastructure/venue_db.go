package infrastructure

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sportradar-be-exercise/services"
)

type VenueRepository struct {
		db *sqlx.DB
}

func NewVenueRepository(db *sqlx.DB) *VenueRepository {
	return &VenueRepository{db: db}
}

func (v *VenueRepository) CreateVenue(ctx context.Context, params services.VenueRequest) (int, error) {
	query := "INSERT INTO venues (name, city, country_code) VALUES ($1, $2, $3) RETURNING id"
	var newID int

	if err := v.db.QueryRowContext(ctx, query, params.Name, params.City,
		 params.CountryCode).Scan(&newID); err != nil {
			return 0, err
	}
	return newID, nil
}

func (v *VenueRepository) GetVenueById(ctx context.Context, id int) (*services.Venue, error) {
	query := "SELECT id, name, city, country_code FROM venues WHERE id = $1"
	var dbModel venueDBModel

	if err := v.db.GetContext(ctx, &dbModel, query, id); err != nil {
		return nil, err
	}
	venue := toServiceVenue(dbModel)
	return &venue, nil
}

func (v *VenueRepository) ListVenues(ctx context.Context) ([]services.Venue, error) {
	query := "SELECT id, name, city, country_code FROM venues ORDER BY name ASC"
	var dbModel []venueDBModel
	
	if err := v.db.SelectContext(ctx, &dbModel, query); err != nil {
		return nil, err
	}
	venues := make([]services.Venue, 0, len(dbModel))
	for _, dbVenue := range dbModel {
		venues = append(venues, toServiceVenue(dbVenue))
	}
	return venues, nil
}

func (v *VenueRepository) UpdateVenue(ctx context.Context, venue services.Venue) error {
	query := "UPDATE venues SET name = $1, city = $2, country_code = $3 WHERE id = $4"

	_, err := v.db.ExecContext(ctx, query, venue.Name, venue.City, venue.CountryCode, venue.ID)
	return err
}

func (v *VenueRepository) DeleteVenue(ctx context.Context, id int) error {
	query := "DELETE FROM venues WHERE id = $1"

	_, err := v.db.ExecContext(ctx, query, id)
	return err
}
