package infrastructure

import (
	"context"
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/vsennikov/sports-event-calendar/config"
)

// SetupTestDB creates a test database connection
func SetupTestDB(t *testing.T) *sqlx.DB {
	t.Helper()

	cfg := config.Config{
		DBHost:     getEnvOrDefault("TEST_DB_HOST", "localhost"),
		DBPort:     getEnvOrDefault("TEST_DB_PORT", "5432"),
		DBUser:     getEnvOrDefault("TEST_DB_USER", "postgres"),
		DBPassword: getEnvOrDefault("TEST_DB_PASSWORD", "postgres"),
		DBName:     getEnvOrDefault("TEST_DB_NAME", "sports_event_calendar_test"),
	}

	db, err := NewConnection(cfg)
	if err != nil {
		t.Skipf("Skipping integration test - failed to connect to test database: %v", err)
		return nil
	}

	return db
}

// CleanupTestDB cleans up test data
func CleanupTestDB(t *testing.T, db *sqlx.DB) {
	t.Helper()

	ctx := context.Background()

	// Delete in reverse order of dependencies
	_, err := db.ExecContext(ctx, "DELETE FROM events")
	if err != nil {
		t.Logf("Error cleaning up events: %v", err)
	}

	_, err = db.ExecContext(ctx, "DELETE FROM teams")
	if err != nil {
		t.Logf("Error cleaning up teams: %v", err)
	}

	_, err = db.ExecContext(ctx, "DELETE FROM venues")
	if err != nil {
		t.Logf("Error cleaning up venues: %v", err)
	}

	_, err = db.ExecContext(ctx, "DELETE FROM sports")
	if err != nil {
		t.Logf("Error cleaning up sports: %v", err)
	}

	// Reset sequences
	_, err = db.ExecContext(ctx, "ALTER SEQUENCE events_id_seq RESTART WITH 1")
	if err != nil {
		t.Logf("Error resetting events sequence: %v", err)
	}

	_, err = db.ExecContext(ctx, "ALTER SEQUENCE teams_id_seq RESTART WITH 1")
	if err != nil {
		t.Logf("Error resetting teams sequence: %v", err)
	}

	_, err = db.ExecContext(ctx, "ALTER SEQUENCE venues_id_seq RESTART WITH 1")
	if err != nil {
		t.Logf("Error resetting venues sequence: %v", err)
	}

	_, err = db.ExecContext(ctx, "ALTER SEQUENCE sports_id_seq RESTART WITH 1")
	if err != nil {
		t.Logf("Error resetting sports sequence: %v", err)
	}
}

// InitTestSchema initializes the test database schema
func InitTestSchema(t *testing.T, db *sqlx.DB) {
	t.Helper()

	ctx := context.Background()

	schema := `
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
	`

	_, err := db.ExecContext(ctx, schema)
	if err != nil {
		t.Fatalf("Failed to initialize test schema: %v", err)
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
