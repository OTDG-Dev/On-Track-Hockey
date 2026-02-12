package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/stats"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

type Player struct {
	ID            int  `json:"id"`
	IsActive      bool `json:"is_active"`
	CurrentTeamId int  `json:"current_team_id"`

	FirstName     string               `json:"first_name"`
	LastName      string               `json:"last_name"`
	SweaterNumber uint8                `json:"sweater_number"`
	Position      Position             `json:"position"`
	BirthDate     BirthDate            `json:"birth_date"`
	BirthCountry  string               `json:"birth_country"`
	Headshot      string               `json:"headshot,omitzero"`
	ShootsCatches ShootsCatches        `json:"shoots_catches,omitzero"`
	SkaterStats   *stats.SkaterStatSet `json:"skater_stats,omitzero"`
	GoalieStats   *stats.GoalieStatSet `json:"goalie_stats,omitzero"`
}

func ValidatePlayer(v *validator.Validator, player *Player) {
	v.Check(player.FirstName != "", "first_name", "must be provided")
	v.Check(player.LastName != "", "last_name", "must be provided")

	v.Check(player.SweaterNumber >= 1, "sweater_number", "must be greater than 0")
	v.Check(player.SweaterNumber <= 100, "sweater_number", "must be less than 100")

	v.Check(player.BirthDate.Year() <= time.Now().Year(), "birth_year", "cannot be in the future")

	v.Check(len(player.BirthCountry) <= 3, "birth_country", "must only be 3 chars")
}

// wrap a sql.DB connection pool
type PlayerModel struct {
	DB *sql.DB
}

func (m PlayerModel) Insert(player *Player) error {
	query := `
		INSERT INTO players (
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
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id`

	args := []any{
		player.IsActive,
		player.CurrentTeamId,
		player.FirstName,
		player.LastName,
		player.SweaterNumber,
		player.Position,
		player.BirthDate,
		player.BirthCountry,
		player.Headshot,
		player.ShootsCatches,
	}

	return m.DB.QueryRow(query, args...).Scan(&player.ID)
}

func (m PlayerModel) Get(id int) (*Player, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
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
		WHERE id = $1`

	var player Player

	err := m.DB.QueryRow(query, id).Scan(
		&player.ID,
		&player.IsActive,
		&player.CurrentTeamId,
		&player.FirstName,
		&player.LastName,
		&player.SweaterNumber,
		&player.Position,
		&player.BirthDate,
		&player.BirthCountry,
		&player.Headshot,
		&player.ShootsCatches,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &player, nil
}

func (m PlayerModel) GetAll() ([]Player, error) {
	// WIP need to catch err no rows exception!
	query := `
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
	FROM players`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	var players []Player

	for rows.Next() {
		var p Player
		err = rows.Scan(
			&p.ID,
			&p.IsActive,
			&p.CurrentTeamId,
			&p.FirstName,
			&p.LastName,
			&p.SweaterNumber,
			&p.Position,
			&p.BirthDate,
			&p.BirthCountry,
			&p.Headshot,
			&p.ShootsCatches,
		)
		players = append(players, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return players, err
}

func (m PlayerModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM players
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
		return err
	}

	return nil
}
