package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
	"github.com/lib/pq"
)

type Division struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"version"`

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

func (m DivisonModel) Get(id int) (*Division, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := /* sql */ `
		SELECT id, league_id, name, version
		FROM divisions
		WHERE id = $1`

	var d Division

	err := m.DB.QueryRow(query, id).Scan(
		&d.ID,
		&d.LeagueID,
		&d.Name,
		&d.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &d, nil
}

func (m DivisonModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := /* sql */ `
		DELETE FROM divisions
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
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

func (m DivisonModel) Update(division *Division) error {
	query := /* sql */ `
		UPDATE divisions
		SET 
			name = $1,
			league_id = $2,
			version = version + 1
		WHERE id = $3 AND version = $4
		RETURNING version`

	args := []any{division.Name, division.LeagueID, division.ID, division.Version}

	err := m.DB.QueryRow(query, args...).Scan(&division.Version)
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
