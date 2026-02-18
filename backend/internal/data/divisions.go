package data

import (
	"database/sql"
	"errors"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
	"github.com/lib/pq"
)

type Division struct {
	ID       int    `json:"id"`
	LeagueID int    `json:"league_id"`
	Name     string `json:"name"`
}

type DivisonModel struct {
	DB *sql.DB
}

func ValidateDivision(v *validator.Validator, div *Division) {
	v.Check(div.LeagueID != 0, "league_id", "league_id must be provided")
}

func (m DivisonModel) Insert(div *Division) error {
	query := /* sql */ `
		INSERT INTO divisions (
			name,
			league_id
		)
		VALUES ($1, $2)
		RETURNING ID`

	err := m.DB.QueryRow(query, []any{div.Name, div.LeagueID}...).Scan(&div.ID)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) {
			if pqErr.Code.Name() == "foreign_key_violation" {
				return ErrNotFound
			}
		}
		return err
	}

	return nil
}

func (m DivisonModel) GetAll() ([]*Division, error) {
	// should query league_id for divs
	query := /* sql */ `
		SELECT
			id,
			name,
			league_id
		FROM divisions`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	divs := []*Division{}

	for rows.Next() {
		var d Division
		err := rows.Scan(&d.ID, &d.Name, &d.LeagueID)
		if err != nil {
			return nil, err
		}
		divs = append(divs, &d)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return divs, nil
}
