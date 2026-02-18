package data

import (
	"database/sql"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

type League struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LeagueModel struct {
	DB *sql.DB
}

func ValidateLeague(v *validator.Validator, league *League) {
	// no rules yet
}

func (m LeagueModel) Insert(league *League) error {
	query := /* sql */ `
		INSERT INTO leagues (
			name
		)
		VALUES ($1)
		RETURNING ID`

	return m.DB.QueryRow(query, []any{league.Name}...).Scan(&league.ID)
}

func (m LeagueModel) GetAll() ([]*League, error) {
	query := /* sql */ `SELECT id, name FROM leagues`

	rows, err := m.DB.Query(query)
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
