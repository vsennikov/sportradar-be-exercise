SET TIMEZONE='UTC';

CREATE TABLE IF NOT EXISTS sports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS venues (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    city VARCHAR(100) NOT NULL,
    country_code CHAR(2) NOT NULL
);

CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    city VARCHAR(100) NOT NULL,
    _sport_id INTEGER NOT NULL,

    CONSTRAINT fk_sport FOREIGN KEY(_sport_id) REFERENCES sports(id),

    CONSTRAINT uq_team_sport UNIQUE (name, _sport_id)
);

CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    event_datetime TIMESTAMPTZ NOT NULL,
    description TEXT,
    home_score INTEGER,
    away_score INTEGER,
    _sport_id INTEGER NOT NULL,
    _venue_id INTEGER,
    _home_team_id INTEGER NOT NULL,
    _away_team_id INTEGER NOT NULL,
    
    CONSTRAINT fk_sport FOREIGN KEY(_sport_id) REFERENCES sports(id),
    CONSTRAINT fk_venue FOREIGN KEY(_venue_id) REFERENCES venues(id),
    CONSTRAINT fk_home_team FOREIGN KEY(_home_team_id) REFERENCES teams(id),
    CONSTRAINT fk_away_team FOREIGN KEY(_away_team_id) REFERENCES teams(id),
    
    CONSTRAINT check_teams_not_equal CHECK (_home_team_id <> _away_team_id)
);


INSERT INTO sports (name) VALUES ('Football'), ('Ice Hockey')
ON CONFLICT (name) DO NOTHING;

INSERT INTO venues (name, city, country_code) VALUES 
('Red Bull Arena', 'Salzburg', 'AT'),
('Etihad Stadium', 'Manchester', 'GB'),
('Parc des Princes', 'Paris', 'FR')
ON CONFLICT DO NOTHING;

INSERT INTO venues (name, city, country_code) VALUES 
('Steffl Arena', 'Vienna', 'AT'),
('Swiss Life Arena', 'Zurich', 'CH'),
('Mercedes-Benz Arena', 'Berlin', 'DE')
ON CONFLICT DO NOTHING;


INSERT INTO teams (name, city, _sport_id) VALUES 
('Red Bull Salzburg', 'Salzburg', 1),
('Manchester City', 'Manchester', 1),
('Paris Saint-Germain', 'Paris', 1)
ON CONFLICT (name, _sport_id) DO NOTHING;


INSERT INTO teams (name, city, _sport_id) VALUES 
('Vienna Capitals', 'Vienna', 2),
('ZSC Lions', 'Zurich', 2),
('Eisbären Berlin', 'Berlin', 2)
ON CONFLICT (name, _sport_id) DO NOTHING;

INSERT INTO events (event_datetime, _sport_id, _home_team_id, _away_team_id, _venue_id, home_score, away_score, description) VALUES
('2025-10-01 19:00:00 UTC', 1, 1, 2, 1, 2, 2, 'Champions League, group stage.');

INSERT INTO events (event_datetime, _sport_id, _home_team_id, _away_team_id, _venue_id) VALUES
('2025-12-10 20:00:00 UTC', 1, 3, 1, 3),
('2025-12-15 17:30:00 UTC', 1, 2, 3, 2);


INSERT INTO events (event_datetime, _sport_id, _home_team_id, _away_team_id, _venue_id, home_score, away_score) VALUES
('2025-10-05 18:00:00 UTC', 2, 4, 5, 4, 5, 3);

INSERT INTO events (event_datetime, _sport_id, _home_team_id, _away_team_id, _venue_id) VALUES
('2025-12-12 19:30:00 UTC', 2, 5, 6, 5),
('2025-12-18 20:15:00 UTC', 2, 6, 4, 6);

-- 47 more events to test pagination
INSERT INTO events (
    event_datetime,
    _sport_id,
    _home_team_id,
    _away_team_id,
    _venue_id,
    description
)
SELECT
    '2025-12-20 18:00:00 UTC'::timestamptz + (i || ' days')::interval,
    1,
    (i % 3) + 1,
    ((i + 1) % 3) + 1,
    (i % 3) + 1,
    'Generated match day ' || i
FROM
    generate_series(1, 47) AS i;