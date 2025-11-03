package controllers

import "github.com/vsennikov/sportradar-be-exercise/services"

func toDTOEvent(event services.Event) EventDTO{
	var venue *venueDTO

	if event.Venue.ID != 0 {
		venue = &venueDTO{
			ID: event.Venue.ID,
			Name: event.Venue.Name,
			City: event.Venue.City,
			CountryCode: event.Venue.CountryCode,
		}
	}
	return EventDTO{
		ID: event.ID,
		EventDatetime: event.EventDatetime,
		Description: event.Description,
		HomeScore: event.HomeScore,
		AwayScore: event.AwayScore,
		
		Sport: sportDTO{
			ID: event.Sport.ID,
			Name: event.Sport.Name,
		},

		Venue: venue,

		HomeTeam: teamDTO{
			ID: event.HomeTeam.ID,
			Name: event.HomeTeam.Name,
			City: event.HomeTeam.City,
		},

		AwayTeam: teamDTO{
			ID: event.AwayTeam.ID,
			Name: event.AwayTeam.Name,
			City: event.AwayTeam.City,
		},
	}
}

func toDTOSport(sport services.Sport) sportDTO {
	return sportDTO{
		ID: sport.ID,
		Name: sport.Name,
	}
}