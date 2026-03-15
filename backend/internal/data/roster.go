package data

import (
	"context"
	"database/sql"
	"time"
)

type Roster struct {
	Forwards   []Player `json:"forwards"`
	Defensemen []Player `json:"defensemen"`
	Goalies    []Player `json:"goalies"`
}

type RosterModel struct {
	DB *sql.DB
}

func (m *RosterModel) Get(teamID int) (*Roster, error) {
	query := /* sql */ `
		SELECT 
			id,
			is_active,
			current_team_id,
			first_name,
			last_name,
			sweater_number,
			position,
			birth_date,
			birth_country,
			headshot,
			shoots_catches
		FROM players
		WHERE current_team_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, teamID)
	if err != nil {
		return nil, err
	}

	var r Roster

	for rows.Next() {
		var p Player
		err := rows.Scan(
			&p.ID,
			&p.IsActive,
			&p.CurrentTeamID,
			&p.FirstName,
			&p.LastName,
			&p.SweaterNumber,
			&p.Position,
			&p.BirthDate,
			&p.BirthCountry,
			&p.Headshot,
			&p.ShootsCatches,
		)
		if err != nil {
			return nil, err
		}
		switch p.Position {
		case PositionC, PositionLW, PositionRW:
			r.Forwards = append(r.Forwards, p)
		case PositionD:
			r.Defensemen = append(r.Defensemen, p)
		case PositionG:
			r.Goalies = append(r.Goalies, p)
		}
	}

	return &r, nil
}
