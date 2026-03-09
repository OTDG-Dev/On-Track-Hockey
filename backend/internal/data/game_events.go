package data

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

type GameEvent struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"-"`

	GameID    int    `json:"game_id"`
	Period    int    `json:"period"`
	Clock     int    `json:"clock"`
	EventType string `json:"event_type"`
	TeamID    int    `json:"team_id"`
}

type GameEventModel struct {
	DB *sql.DB
}

var eventTypes = []string{"goal", "penalty", "shot", "save"} // could be a custom type..

func ValidateGameEvent(v *validator.Validator, e *GameEvent) {
	e.EventType = strings.ToLower(e.EventType)
	v.Check(slices.Contains(eventTypes, e.EventType), "event_type", "this type not permitted")
}

func (m *GameEventModel) Insert(event GameEvent) error {
	query := /* sql */ `
		INSERT INTO game_events (
			game_id,
			period,
			clock,
			event_type,
			team_id,
		)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{event.GameID, event.Period, event.Clock, event.EventType, event.TeamID}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&event.ID, &event.CreatedAt, &event.Version)
}

func (m *GameEventModel) Get(id int) (*GameEvent, error) {
	query := /* sql */ `
		SELECT
			version,
			game_id,
			period,
			clock,
			event_type,
			team_id,
		FROM game_events
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event GameEvent

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&event.Version,
		&event.GameID,
		&event.Period,
		&event.Clock,
		&event.EventType,
		&event.TeamID,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &event, err
}
