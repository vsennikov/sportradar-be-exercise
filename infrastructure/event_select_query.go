package infrastructure

const baseEventSelectQuery = `
SELECT
    e.id,
    e.event_datetime,
    e.description,
    e.home_score,
    e.away_score,
    s.id AS "sport.id",
    s.name AS "sport.name",
    v.id AS "venue.id",
    v.name AS "venue.name",
    v.city AS "venue.city",
    v.country_code AS "venue.country_code",
    ht.id AS "ht.id",
    ht.name AS "ht.name",
    ht.city AS "ht.city",
    at.id AS "at.id",
    at.name AS "at.name",
    at.city AS "at.city"
FROM events e
JOIN sports s ON e._sport_id = s.id
LEFT JOIN venues v ON e._venue_id = v.id
JOIN teams ht ON e._home_team_id = ht.id
JOIN teams at ON e._away_team_id = at.id
`