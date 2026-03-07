package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

// EventParticipant represents a player's participation in a game event
type EventParticipant struct {
	ID          int       `json:"id"`
	CreatedAt   time.Time `json:"-"`
	GameEventID int       `json:"game_event_id"`
	PlayerID    int       `json:"player_id"`
	Role        string    `json:"role"`
}

// EventParticipantWithDetails includes player information along with the participation
type EventParticipantWithDetails struct {
	EventParticipant
	PlayerFirstName string `json:"player_first_name"`
	PlayerLastName  string `json:"player_last_name"`
	PlayerNumber    int    `json:"player_number,omitempty"`
}

// EventParticipantModel wraps the database connection for event participant operations
type EventParticipantModel struct {
	DB *sql.DB
}

// Insert creates a new event participant in the database
func (m *EventParticipantModel) Insert(participant *EventParticipant) error {
	query := `
		INSERT INTO event_participants (game_event_id, player_id, role)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{participant.GameEventID, participant.PlayerID, participant.Role}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&participant.ID, &participant.CreatedAt)
}

// Get retrieves a single event participant by ID
func (m *EventParticipantModel) Get(id int) (*EventParticipant, error) {
	query := `
		SELECT id, game_event_id, player_id, role, created_at
		FROM event_participants
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var participant EventParticipant
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&participant.ID,
		&participant.GameEventID,
		&participant.PlayerID,
		&participant.Role,
		&participant.CreatedAt,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &participant, nil
}

// GetWithDetails retrieves an event participant with player information
func (m *EventParticipantModel) GetWithDetails(id int) (*EventParticipantWithDetails, error) {
	query := `
		SELECT 
			ep.id, ep.game_event_id, ep.player_id, ep.role, ep.created_at,
			p.first_name, p.last_name, p.number
		FROM event_participants ep
		INNER JOIN players p ON ep.player_id = p.id
		WHERE ep.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var participant EventParticipantWithDetails
	var playerNumber sql.NullInt32
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&participant.ID,
		&participant.GameEventID,
		&participant.PlayerID,
		&participant.Role,
		&participant.CreatedAt,
		&participant.PlayerFirstName,
		&participant.PlayerLastName,
		&playerNumber,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	if playerNumber.Valid {
		participant.PlayerNumber = int(playerNumber.Int32)
	}

	return &participant, nil
}

// ListByGameEvent retrieves all participants for a specific game event
func (m *EventParticipantModel) ListByGameEvent(gameEventID int) ([]*EventParticipantWithDetails, error) {
	query := `
		SELECT 
			ep.id, ep.game_event_id, ep.player_id, ep.role, ep.created_at,
			p.first_name, p.last_name, p.number
		FROM event_participants ep
		INNER JOIN players p ON ep.player_id = p.id
		WHERE ep.game_event_id = $1
		ORDER BY ep.id ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, gameEventID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*EventParticipantWithDetails
	for rows.Next() {
		var participant EventParticipantWithDetails
		var playerNumber sql.NullInt32
		err := rows.Scan(
			&participant.ID,
			&participant.GameEventID,
			&participant.PlayerID,
			&participant.Role,
			&participant.CreatedAt,
			&participant.PlayerFirstName,
			&participant.PlayerLastName,
			&playerNumber,
		)
		if err != nil {
			return nil, err
		}
		if playerNumber.Valid {
			participant.PlayerNumber = int(playerNumber.Int32)
		}
		participants = append(participants, &participant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}

// ListByPlayer retrieves all event participants for a specific player
func (m *EventParticipantModel) ListByPlayer(playerID int) ([]*EventParticipantWithDetails, error) {
	query := `
		SELECT 
			ep.id, ep.game_event_id, ep.player_id, ep.role, ep.created_at,
			p.first_name, p.last_name, p.number
		FROM event_participants ep
		INNER JOIN players p ON ep.player_id = p.id
		WHERE ep.player_id = $1
		ORDER BY ep.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*EventParticipantWithDetails
	for rows.Next() {
		var participant EventParticipantWithDetails
		var playerNumber sql.NullInt32
		err := rows.Scan(
			&participant.ID,
			&participant.GameEventID,
			&participant.PlayerID,
			&participant.Role,
			&participant.CreatedAt,
			&participant.PlayerFirstName,
			&participant.PlayerLastName,
			&playerNumber,
		)
		if err != nil {
			return nil, err
		}
		if playerNumber.Valid {
			participant.PlayerNumber = int(playerNumber.Int32)
		}
		participants = append(participants, &participant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return participants, nil
}

// Update modifies an existing event participant
func (m *EventParticipantModel) Update(participant *EventParticipant) error {
	query := `
		UPDATE event_participants
		SET player_id = $1, role = $2
		WHERE id = $3`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{participant.PlayerID, participant.Role, participant.ID}
	result, err := m.DB.ExecContext(ctx, query, args...)
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

// Delete removes an event participant by ID
func (m *EventParticipantModel) Delete(id int) error {
	query := `
		DELETE FROM event_participants
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

// DeleteByGameEvent removes all participants for a specific game event
func (m *EventParticipantModel) DeleteByGameEvent(gameEventID int) error {
	query := `
		DELETE FROM event_participants
		WHERE game_event_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.DB.ExecContext(ctx, query, gameEventID)
	return err
}
