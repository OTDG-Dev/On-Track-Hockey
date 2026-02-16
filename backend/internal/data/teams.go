package data

import (
	"database/sql"
	"errors"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

type Team struct {
	ID         int    `json:"id"`
	FullName   string `json:"full_name"`
	ShortName  string `json:"short_name"`
	DivisionID int    `json:"division_id"`
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
		SELECT id, full_name, short_name, division_id
		FROM teams
		WHERE id = $1;`

	var t Team

	err := m.DB.QueryRow(query, id).Scan(
		&t.ID,
		&t.FullName,
		&t.ShortName,
		&t.DivisionID,
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
	query := /* sql */ `
	INSERT INTO teams (
		full_name,
		short_name,
		division_id
	)
	VALUES ($1, $2, $3)
	RETURNING ID`

	args := []any{team.FullName, team.ShortName, team.DivisionID}

	return m.DB.QueryRow(query, args...).Scan(&team.ID)
}

func (m TeamModel) Delete(id int) error {
	query := /* sql */ `
		DELETE FROM teams
		WHERE id = $1
	`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return err
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
			division_id
		FROM teams;`

	rows, err := m.DB.Query(query)
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
