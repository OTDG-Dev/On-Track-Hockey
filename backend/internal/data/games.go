package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type Game struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"-"`

	HomeTeamID int       `json:"home_team_id"`
	AwayTeamID int       `json:"away_team_id"`
	StartTime  time.Time `json:"start_time"`
	IsFinished bool      `json:"is_finished"`
}

type GameModel struct {
	DB *sql.DB
}

func (m *GameModel) Insert(game *Game) error {
	query := /* sql */ `
		INSERT INTO games (
			home_team_id,
			away_team_id,
			start_time,
			is_finished
		)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	args := []any{game.HomeTeamID, game.AwayTeamID, game.StartTime, game.IsFinished}
	return m.DB.QueryRowContext(ctx, query, args...).Scan(&game.ID, &game.CreatedAt, &game.Version)
}

func (m *GameModel) Get(id int) (*Game, error) {
	query := /* sql */ `
		SELECT
			id,
			created_at,
			version,
			home_team_id,
			away_team_id,
			start_time,
			is_finished
		FROM games
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var game Game

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&game.ID,
		&game.CreatedAt,
		&game.Version,
		&game.HomeTeamID,
		&game.AwayTeamID,
		&game.StartTime,
		&game.IsFinished,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &game, nil
}

func (m *GameModel) Update(game *Game) error {
	query := /* sql */ `
		UPDATE games
		SET
			home_team_id = $1,
			away_team_id = $2,
			is_finished = $3,
			version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING version`

	args := []any{
		game.HomeTeamID,
		game.AwayTeamID,
		game.IsFinished,
		game.ID,
		game.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&game.Version)
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

type GameView struct {
	HomeTeam   string      `json:"home_team"`
	AwayTeam   string      `json:"away_team"`
	HomeTeamID int         `json:"home_team_id"`
	AwayTeamID int         `json:"away_team_id"`
	IsFinished bool        `json:"is_finished"`
	StartTime  time.Time   `json:"start_time"`
	GameEvents []GameEvent `json:"game_events"`
}

type GameListView struct {
	ID         int       `json:"id"`
	HomeTeam   string    `json:"home_team"`
	AwayTeam   string    `json:"away_team"`
	HomeTeamID int       `json:"home_team_id"`
	AwayTeamID int       `json:"away_team_id"`
	IsFinished bool      `json:"is_finished"`
	StartTime  time.Time `json:"start_time"`
}

func (m *GameModel) GetView(gameID int) (*GameView, error) {

	// Phase 1: Game Info
	gameQuery := /* sql */ `
		SELECT
			g.home_team_id,
			t1.short_name,
			g.away_team_id,
			t2.short_name,
			g.is_finished,
			g.start_time
		FROM games g
		INNER JOIN teams t1
			ON g.home_team_id = t1.id
		INNER JOIN teams t2
			ON g.away_team_id = t2.id
		WHERE g.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var g GameView

	err := m.DB.QueryRowContext(ctx, gameQuery, gameID).Scan(
		&g.HomeTeamID,
		&g.HomeTeam,
		&g.AwayTeamID,
		&g.AwayTeam,
		&g.IsFinished,
		&g.StartTime,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Phase 2: Game events
	eventsQuery := /* sql */ `
		SELECT
			id,
			period,
			clock_seconds,
			event_type,
			situation,
			event_number,
			team_id
		FROM game_events
		WHERE game_id = $1`

	ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, eventsQuery, gameID)
	if err != nil {
		return &g, nil
	}
	defer rows.Close()

	gameEvents := []GameEvent{}

	for rows.Next() {
		var e GameEvent
		err := rows.Scan(
			&e.ID,
			&e.Period,
			&e.ClockSeconds,
			&e.EventType,
			&e.Situation,
			&e.EventNumber,
			&e.TeamID,
		)
		if err != nil {
			return &g, err
		}
		gameEvents = append(gameEvents, e)
	}

	if err := rows.Err(); err != nil {
		return &g, err
	}

	g.GameEvents = gameEvents

	return &g, nil
}

func (m *GameModel) GetAll() ([]*GameListView, error) {
	query := /* sql */ `
		SELECT
			g.id,
			g.home_team_id,
			t1.short_name,
			g.away_team_id,
			t2.short_name,
			g.is_finished,
			g.start_time
		FROM games g
		INNER JOIN teams t1
			ON g.home_team_id = t1.id
		INNER JOIN teams t2
			ON g.away_team_id = t2.id
		ORDER BY g.start_time DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	games := []*GameListView{}

	for rows.Next() {
		var g GameListView
		err = rows.Scan(
			&g.ID,
			&g.HomeTeamID,
			&g.HomeTeam,
			&g.AwayTeamID,
			&g.AwayTeam,
			&g.IsFinished,
			&g.StartTime,
		)
		if err != nil {
			return nil, err
		}
		games = append(games, &g)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return games, nil
}
