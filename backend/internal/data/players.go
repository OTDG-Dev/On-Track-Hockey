package data

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/stats"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

type Player struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"version"`

	IsActive      bool `json:"is_active"`
	CurrentTeamID *int `json:"current_team_id"`

	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	SweaterNumber uint8         `json:"sweater_number"`
	Position      Position      `json:"position"`
	BirthDate     BirthDate     `json:"birth_date"`
	BirthCountry  string        `json:"birth_country"`
	Headshot      string        `json:"headshot,omitzero"`
	ShootsCatches ShootsCatches `json:"shoots_catches,omitzero"`

	SkaterStats *stats.SkaterStatSet `json:"skater_stats,omitzero"`
	GoalieStats *stats.GoalieStatSet `json:"goalie_stats,omitzero"`
}

// wrap a sql.DB connection pool
type PlayerModel struct {
	DB *sql.DB
}

func ValidatePlayer(v *validator.Validator, player *Player) {
	v.Check(player.FirstName != "", "first_name", "must be provided")
	v.Check(player.LastName != "", "last_name", "must be provided")

	v.Check(player.SweaterNumber >= 1, "sweater_number", "must be greater than 0")
	v.Check(player.SweaterNumber <= 100, "sweater_number", "must be less than 100")

	v.Check(player.BirthDate.Year() <= time.Now().Year(), "birth_year", "cannot be in the future")

	v.Check(len(player.BirthCountry) <= 3, "birth_country", "must only be 3 chars")

	v.Check(validator.PermittedValue(player.Position, "C", "LW", "RW", "D", "G"), "position", "must be 'C|LW|RW|D|G'")
	v.Check(validator.PermittedValue(player.ShootsCatches, "L", "R"), "shoots_catches", "must be 'L|R'")
}

func (m PlayerModel) Insert(player *Player) error {
	query := /* sql */ `
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
		player.CurrentTeamID,
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
			shoots_catches,
			version
		FROM players
		WHERE id = $1`

	var player Player

	err := m.DB.QueryRow(query, id).Scan(
		&player.ID,
		&player.IsActive,
		&player.CurrentTeamID,
		&player.FirstName,
		&player.LastName,
		&player.SweaterNumber,
		&player.Position,
		&player.BirthDate,
		&player.BirthCountry,
		&player.Headshot,
		&player.ShootsCatches,
		&player.Version,
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

func (m PlayerModel) GetAll(FirstName, LastName, Position string, filters Filters) ([]*Player, Metadata, error) {
	// WIP need to use like and also combine first/lastname into the query
	// https://niallburkley.com/blog/index-columns-for-like-in-postgres/
	query := fmt.Sprintf( /* sql */ `
	SELECT
		count(*) OVER(),
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
		shoots_catches,
		version
	FROM players
	WHERE (first_name ILIKE $1 OR $1 = '')  -- switch to indexes with scale & combine last + first
	AND (last_name ILIKE $2 OR $2 = '')
	AND (position = $3 OR $3 = '')
	ORDER BY %s %s, id ASC
	LIMIT $4 OFFSET $5`, filters.sortColumn(), filters.SortDirection())

	args := []any{FirstName, LastName, Position, filters.limit(), filters.offset()}

	rows, err := m.DB.Query(query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	players := []*Player{}

	for rows.Next() {
		var p Player
		err = rows.Scan(
			&totalRecords,
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
			&p.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}

		players = append(players, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	metadata := calculateMetadata(totalRecords, filters.Page, filters.PageSize)

	return players, metadata, err
}

func (m PlayerModel) Update(player *Player) error {
	query := /* sql */ `
		UPDATE players
		SET
			is_active = $1,
			current_team_id = $2,
			first_name = $3,
			last_name = $4,
			sweater_number = $5,
			position = $6,
			birth_date = $7,
			birth_country = $8,
			headshot = $9,
			shoots_catches = $10,
			version = version + 1
		WHERE id = $11 AND version = $12
		RETURNING version`

	args := []any{
		player.IsActive,
		player.CurrentTeamID,
		player.FirstName,
		player.LastName,
		player.SweaterNumber,
		player.Position,
		player.BirthDate,
		player.BirthCountry,
		player.Headshot,
		player.ShootsCatches,
		player.ID,
		player.Version,
	}

	fmt.Println(player.ID, player.Version)

	err := m.DB.QueryRow(query, args...).Scan(&player.Version)
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

func (m PlayerModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := /* sql */ `
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
