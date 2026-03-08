package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

type League struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"version"`

	Name string `json:"name"`
}

type LeagueModel struct {
	DB *sql.DB
}

func ValidateLeague(v *validator.Validator, league *League) {
	v.Check(len(league.Name) <= 128 && len(league.Name) > 0, "name", "league name must be between 1-128 characters")
}

func (m LeagueModel) Get(id int) (*League, error) {
	query := /* sql */ `
		SELECT id, name, version
		FROM leagues
		WHERE id = $1`

	var l League

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(&l.ID, &l.Name, &l.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &l, nil
}

func (m LeagueModel) Insert(league *League) error {
	query := /* sql */ `
		INSERT INTO leagues (
			name
		)
		VALUES ($1)
		RETURNING ID`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, []any{league.Name}...).Scan(&league.ID)
}

func (m LeagueModel) GetAll() ([]*League, error) {
	query := /* sql */ `SELECT id, name FROM leagues`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	leagues := []*League{}

	for rows.Next() {
		var l League
		err = rows.Scan(&l.ID, &l.Name)
		if err != nil {
			return nil, err
		}
		leagues = append(leagues, &l)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return leagues, nil
}

func (m LeagueModel) Update(league *League) error {
	query := /* sql */ `
		UPDATE leagues
		SET
			name = $1,
			version = version + 1
		WHERE id = $2 AND version = $3
		RETURNING version`

	args := []any{league.Name, league.ID, league.Version}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&league.Version)
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

func (m LeagueModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := /* sql */ `
		DELETE FROM leagues
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return ExecDeleteErrors(err, "league")
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
