package infrastructure

import (
	"database/sql"

	"github.com/vsennikov/sportradar-be-exercise/services"
)

func toServiceEvent(db eventDBModel) services.Event {
	var venue services.Venue
	if db.VenueID.Valid {
		venue = services.Venue{
			ID:          int(db.VenueID.Int64),
			Name:        db.VenueName.String,
			City:        db.VenueCity.String,
			CountryCode: db.VenueCountryCode.String,
		}
	}

	return services.Event{
		ID:            db.ID,
		EventDatetime: db.EventDatetime,
		Description:   nullStringToStringPtr(db.Description),
		HomeScore:     nullInt64ToIntPtr(db.HomeScore),
		AwayScore:     nullInt64ToIntPtr(db.AwayScore),
		Sport: services.Sport{
			ID:   db.SportID,
			Name: db.SportName,
		},
		Venue: venue,
		HomeTeam: services.Team{
			ID:   db.HomeTeamID,
			Name: db.HomeTeamName,
			City: db.HomeTeamCity,
		},
		AwayTeam: services.Team{
			ID:   db.AwayTeamID,
			Name: db.AwayTeamName,
			City: db.AwayTeamCity,
		},
	}
}

func nullStringToStringPtr(s sql.NullString) *string {
	if s.Valid {
		return &s.String
	}
	return nil
}

func nullInt64ToIntPtr(i sql.NullInt64) *int {
	if i.Valid {
		val := int(i.Int64)
		return &val
	}
	return nil
}

func toServiceSport(db sportDBModel) services.Sport {
	return services.Sport{
		ID: db.ID,
		Name: db.Name,
	}
}

func toServiceVenue(db venueDBModel) services.Venue {
	return services.Venue{
		ID: db.ID,
		Name: db.Name,
		City: db.City,
		CountryCode: db.CountryCode,
	}
}

func toServiceTeam(db teamDBModel) services.Team {
	return services.Team{
		ID: db.ID,
		Name: db.Name,
		City: db.City,
	}
}
