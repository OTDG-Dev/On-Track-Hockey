package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// GameEvent represents a single event that occurs during a game
type GameEvent struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	GameID    int       `json:"game_id"`
	Period    int       `json:"period"`
	Clock     string    `json:"clock"`
	EventType string    `json:"event_type"`
	TeamID    int       `json:"team_id"`
	Version   int       `json:"version"`
}

// GameEventWithDetails includes team information along with the event
type GameEventWithDetails struct {
	GameEvent
	TeamName string `json:"team_name"`
}

// GameEventModel wraps the database connection for game event operations
type GameEventModel struct {
	DB *sql.DB
}

// Insert creates a new game event in the database
func (m *GameEventModel) Insert(event *GameEvent) error {
	query := `
		INSERT INTO game_events (game_id, period, clock, event_type, team_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{event.GameID, event.Period, event.Clock, event.EventType, event.TeamID}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&event.ID, &event.CreatedAt, &event.Version)
}

// Get retrieves a single game event by ID
func (m *GameEventModel) Get(id int) (*GameEvent, error) {
	query := `
		SELECT id, game_id, period, clock, event_type, team_id, created_at, version
		FROM game_events
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event GameEvent
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.GameID,
		&event.Period,
		&event.Clock,
		&event.EventType,
		&event.TeamID,
		&event.CreatedAt,
		&event.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &event, nil
}

// GetWithDetails retrieves a game event with team name
func (m *GameEventModel) GetWithDetails(id int) (*GameEventWithDetails, error) {
	query := `
		SELECT 
			ge.id, ge.game_id, ge.period, ge.clock, ge.event_type, ge.team_id, ge.created_at, ge.version,
			t.short_name
		FROM game_events ge
		INNER JOIN teams t ON ge.team_id = t.id
		WHERE ge.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event GameEventWithDetails
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.GameID,
		&event.Period,
		&event.Clock,
		&event.EventType,
		&event.TeamID,
		&event.CreatedAt,
		&event.Version,
		&event.TeamName,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &event, nil
}

// ListByGame retrieves all game events for a specific game
func (m *GameEventModel) ListByGame(gameID int) ([]*GameEventWithDetails, error) {
	query := `
		SELECT 
			ge.id, ge.game_id, ge.period, ge.clock, ge.event_type, ge.team_id, ge.created_at, ge.version,
			t.short_name
		FROM game_events ge
		INNER JOIN teams t ON ge.team_id = t.id
		WHERE ge.game_id = $1
		ORDER BY ge.period ASC, ge.clock ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, gameID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*GameEventWithDetails
	for rows.Next() {
		var event GameEventWithDetails
		err := rows.Scan(
			&event.ID,
			&event.GameID,
			&event.Period,
			&event.Clock,
			&event.EventType,
			&event.TeamID,
			&event.CreatedAt,
			&event.Version,
			&event.TeamName,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Update modifies an existing game event
func (m *GameEventModel) Update(event *GameEvent) error {
	query := `
		UPDATE game_events
		SET period = $1, clock = $2, event_type = $3, team_id = $4, version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{event.Period, event.Clock, event.EventType, event.TeamID, event.ID, event.Version}
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&event.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}

	return nil
}

// Delete removes a game event by ID
func (m *GameEventModel) Delete(id int) error {
	query := `
		DELETE FROM game_events
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
