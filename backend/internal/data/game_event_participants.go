package data

import (
	"context"
	"database/sql"
	"slices"
	"strings"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"
)

type GameEventParticipant struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"-"`

	Role     string `json:"role"`
	EventID  int    `json:"event_id"`
	PlayerID int    `json:"player_id"`
}

type GameEventParticipantModel struct {
	DB *sql.DB
}

var eventRoles = []string{"scorer", "assist_primary", "assist_secondary", "penalty_taker"}

func ValidateGameEventParticiant(v *validator.Validator, part *GameEventParticipant) {
	part.Role = strings.ToLower(part.Role)
	msg := "must be one of:" + strings.Join(eventTypes, ",")
	v.Check(slices.Contains(eventRoles, part.Role), "part", msg)

	v.Check(part.PlayerID != 0, "player_id", "must not be empty")
}

func (m *GameEventParticipantModel) Insert(part GameEventParticipant) error {
	query := /* sql */ `
		INSERT INTO game_event_participants (
			role,
			event_id,
			player_id
		)
		VALUES ($1, $2, $3)
		RETURNING ID`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, part.Role, part.EventID, part.PlayerID).Scan(&part.ID)
}
