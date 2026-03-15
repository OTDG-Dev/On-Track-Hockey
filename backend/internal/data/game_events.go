package data

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"
)

type GameEvent struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"-"`

	EventNumber  int    `json:"event_number"`
	GameID       int    `json:"game_id,omitzero"`
	Period       int    `json:"period"`
	ClockSeconds int    `json:"clock_seconds"`
	EventType    string `json:"event_type"`
	Situation    string `json:"situation"`
	TeamID       int    `json:"team_id"`
}

type GameEventModel struct {
	DB *sql.DB
}

const maxClockSeconds = 1200

var eventTypes = []string{"goal", "penalty", "shot", "save"}
var situations = []string{"EV", "PP", "SH", "EN"}

func ValidateGameEvent(v *validator.Validator, e *GameEvent) {
	e.EventType = strings.ToLower(e.EventType)
	msg := "must be one of:" + strings.Join(eventTypes, ",")
	v.Check(slices.Contains(eventTypes, e.EventType), "event_type", msg)

	e.Situation = strings.ToUpper(e.Situation)
	msg = "must be one of: " + strings.Join(situations, ",")
	v.Check(slices.Contains(situations, e.Situation), "situation", msg)

	v.Check(e.Period >= 1 && e.Period <= 3, "period", "invalid period")

	msg = "invalid clock, must be less than " + strconv.Itoa(maxClockSeconds)
	v.Check(e.ClockSeconds >= 0 && e.ClockSeconds <= maxClockSeconds, "clock_seconds", msg)

	v.Check(e.TeamID > 0, "team_id", "must be greater than 0")
}

func (m *GameEventModel) Insert(event *GameEvent) error {
	query := /* sql */ `
		INSERT INTO game_events (
			game_id,
			event_number,
			period,
			clock_seconds,
			event_type,
			situation,
			team_id
		)
		SELECT
			$1,
			COALESCE(MAX(event_number),0) + 1, -- Numbering name space for local game
			$2,
			$3,
			$4,
			$5,
			$6
		FROM game_events
		WHERE game_id = $1
		RETURNING id, created_at, version, event_number`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{
		event.GameID,
		event.Period,
		event.ClockSeconds,
		event.EventType,
		event.Situation,
		event.TeamID,
	}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(
		&event.ID,
		&event.CreatedAt,
		&event.Version,
		&event.EventNumber,
	)
}

func (m *GameEventModel) Get(id int) (*GameEvent, error) {
	query := /* sql */ `
		SELECT
			id,
			event_number,
			version,
			game_id,
			period,
			clock_seconds,
			event_type,
			situation,
			team_id
		FROM game_events
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var event GameEvent

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.EventNumber,
		&event.Version,
		&event.GameID,
		&event.Period,
		&event.ClockSeconds,
		&event.EventType,
		&event.Situation,
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

	return &event, nil
}
