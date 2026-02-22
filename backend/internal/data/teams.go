package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
	"github.com/lib/pq"
)

type Team struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"version"`

	FullName   string `json:"full_name"`
	ShortName  string `json:"short_name"`
	DivisionID int    `json:"division_id"`
	IsActive   bool   `json:"is_active"`
}

type TeamModel struct {
	DB *sql.DB
}

func ValidateTeam(v *validator.Validator, team *Team) {
	v.Check(len(team.ShortName) == 3, "short_name", "short name must have a length of 3")
	v.Check(team.DivisionID != 0, "division_id", "division_id must be provided")
}

func (m TeamModel) Get(id int) (*Team, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := /* sql */ `
		SELECT id, full_name, short_name, division_id, version
		FROM teams
		WHERE id = $1;`

	var t Team

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&t.ID,
		&t.FullName,
		&t.ShortName,
		&t.DivisionID,
		&t.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &t, nil
}

func (m TeamModel) Insert(team *Team) error {
	fmt.Println(team)

	query := /* sql */ `
	INSERT INTO teams (
		full_name,
		short_name,
		division_id,
		is_active
	)
	VALUES ($1, $2, $3, $4)
	RETURNING ID`

	args := []any{team.FullName, team.ShortName, team.DivisionID, team.IsActive}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&team.ID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "foreign_key_violation" {
				return ErrRecordNotFound
			}
		}
		return err
	}

	return nil
}

func (m TeamModel) Update(team *Team) error {
	query := /* sql */ `
		UPDATE teams
		SET
			full_name = $1,
			short_name = $2,
			division_id = $3,
			is_active = $4,
			version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING version`

	args := []any{
		team.FullName,
		team.ShortName,
		team.DivisionID,
		team.IsActive,
		team.ID,
		team.Version,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&team.Version)
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

func (m TeamModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := /* sql */ `
		DELETE FROM teams
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

func (m TeamModel) GetAll() ([]*Team, error) {
	// add division filter
	query := /* sql */ `
		SELECT
			id,
			full_name,
			short_name,
			division_id,
			is_active
		FROM teams;`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := []*Team{}

	for rows.Next() {
		var t Team
		err = rows.Scan(
			&t.ID,
			&t.FullName,
			&t.ShortName,
			&t.DivisionID,
			&t.IsActive,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// add metadata

	return teams, nil
}
