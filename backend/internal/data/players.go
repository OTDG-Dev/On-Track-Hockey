package data

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/stats"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"
	"github.com/lib/pq"
)

type Player struct {
	ID        int       `json:"id"`
	CreatedAt time.Time `json:"-"`
	Version   int       `json:"-"`

	IsActive      bool `json:"is_active"`
	CurrentTeamID int  `json:"current_team_id"`

	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	SweaterNumber uint8         `json:"sweater_number"`
	Position      Position      `json:"position"`
	BirthDate     BirthDate     `json:"birth_date"`
	BirthCountry  string        `json:"birth_country"`
	Headshot      string        `json:"headshot,omitzero"`
	ShootsCatches ShootsCatches `json:"shoots_catches,omitzero"`

	SkaterStats *stats.SkaterStatSet `json:"skater_stats"`
	GoalieStats *stats.GoalieStatSet `json:"goalie_stats"`
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

	v.Check(player.CurrentTeamID != 0, "current_team_id", "must be provided")
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

	err := m.DB.QueryRow(query, args...).Scan(&player.ID)
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

func (m PlayerModel) Get(id int) (*Player, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

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
			shoots_catches,
			version
		FROM players
		WHERE id = $1`

	var p Player

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &p, nil
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

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&player.Version)
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
		return ErrRecordNotFound
	}

	return nil
}

type PlayerWithTeam struct {
	Player
	TeamFullName  string `json:"team_full_name"`
	TeamShortName string `json:"team_short_name"`
}

type PlayerQuery struct {
	FirstName     string
	LastName      string
	Position      string
	CurrentTeamID int
}

type playerStatTotals struct {
	PlayerID    int
	GamesPlayed int
	Goals       int
	Assists     int
	PIM         int
}

func (m PlayerModel) GetView(id int) (*PlayerWithTeam, error) {
	query := /* sql */ `
		SELECT
			p.id,
			p.is_active,
			p.current_team_id,
			p.first_name,
			p.last_name,
			p.sweater_number,
			p.position,
			p.birth_date,
			p.birth_country,
			p.headshot,
			p.shoots_catches,
			p.version,
			t.full_name,
			t.short_name,
			ps.skater_stats::text,
			ps.goalie_stats::text
		FROM players p
		INNER JOIN teams t
			ON p.current_team_id = t.id
		LEFT JOIN player_stats ps
			ON p.id = ps.player_id
		WHERE p.id = $1
	`

	var p PlayerWithTeam
	var skaterStatsJSON sql.NullString
	var goalieStatsJSON sql.NullString

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
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
		&p.TeamFullName,
		&p.TeamShortName,
		&skaterStatsJSON,
		&goalieStatsJSON,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	// Load saved skater stats
	if skaterStatsJSON.Valid && skaterStatsJSON.String != "" {
		p.SkaterStats = &stats.SkaterStatSet{}
		if err := json.Unmarshal([]byte(skaterStatsJSON.String), p.SkaterStats); err != nil {
			return nil, err
		}
	}

	// Load saved goalie stats
	if goalieStatsJSON.Valid && goalieStatsJSON.String != "" {
		p.GoalieStats = &stats.GoalieStatSet{}
		if err := json.Unmarshal([]byte(goalieStatsJSON.String), p.GoalieStats); err != nil {
			return nil, err
		}
	}

	// Default to zero-value stats when none exist
	if p.Position != "G" && p.SkaterStats == nil {
		p.SkaterStats = &stats.SkaterStatSet{}
	}
	if p.Position == "G" && p.GoalieStats == nil {
		p.GoalieStats = &stats.GoalieStatSet{}
	}

	return &p, nil
}

// RebuildStats will recalculate and update all player stats based on game events. It returns the number of players updated
func (m PlayerModel) RebuildStats() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Phase 1: Start rebuild transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	// Phase 2: Clear old stats
	_, err = tx.ExecContext(ctx, `DELETE FROM player_stats`)
	if err != nil {
		return 0, err
	}

	// Phase 3: Read player totals
	query := /* sql */ `
		SELECT
			p.id,
			COUNT(DISTINCT ge.game_id),
			COUNT(*) FILTER (WHERE ge.event_type = 'goal' AND gep.role = 'scorer'),
			COUNT(*) FILTER (WHERE ge.event_type = 'goal' AND gep.role IN ('assist_primary', 'assist_secondary')),
			COUNT(*) FILTER (WHERE ge.event_type = 'penalty' AND gep.role = 'penalty_taker')
		FROM players p
		LEFT JOIN game_event_participants gep ON gep.player_id = p.id
		LEFT JOIN game_events ge ON ge.id = gep.event_id
		WHERE p.position <> 'G'
		GROUP BY p.id`

	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return 0, err
	}

	insertQuery := /* sql */ `
		INSERT INTO player_stats (player_id, skater_stats)
		VALUES ($1, $2)`

	allTotals := []playerStatTotals{}

	for rows.Next() {
		var totals playerStatTotals

		err = rows.Scan(
			&totals.PlayerID,
			&totals.GamesPlayed,
			&totals.Goals,
			&totals.Assists,
			&totals.PIM,
		)
		if err != nil {
			rows.Close()
			return 0, err
		}

		allTotals = append(allTotals, totals)
	}

	if err = rows.Err(); err != nil {
		rows.Close()
		return 0, err
	}

	if err = rows.Close(); err != nil {
		return 0, err
	}

	// Phase 4: Build and save stats
	updated := 0

	for _, totals := range allTotals {
		if totals.GamesPlayed == 0 && totals.Goals == 0 && totals.Assists == 0 && totals.PIM == 0 {
			continue
		}

		basicStats := stats.SkaterStats{
			GamesPlayed: totals.GamesPlayed,
			Goals:       totals.Goals,
			Assists:     totals.Assists,
			Points:      totals.Goals + totals.Assists,
			PIM:         totals.PIM,
		}

		statSet := stats.SkaterStatSet{
			CurrentSeason: basicStats,
			CareerTotals:  basicStats,
		}

		payload, err := json.Marshal(statSet)
		if err != nil {
			return 0, err
		}

		_, err = tx.ExecContext(ctx, insertQuery, totals.PlayerID, payload)
		if err != nil {
			return 0, err
		}

		updated++
	}

	// Phase 5: Commit rebuild
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return updated, nil
}

func (m PlayerModel) GetViewAll(pq PlayerQuery, filters Filters) ([]*PlayerWithTeam, Metadata, error) {
	query := fmt.Sprintf( /* sql */ `
		SELECT
			count(*) OVER(),
			p.id,
			p.is_active,
			p.current_team_id,
			p.first_name,
			p.last_name,
			p.sweater_number,
			p.position,
			p.birth_date,
			p.birth_country,
			p.headshot,
			p.shoots_catches,
			p.version,
			t.full_name,
			t.short_name
		FROM players p
		INNER JOIN teams t
			ON p.current_team_id = t.id
		WHERE (first_name ILIKE $1 OR $1 = '')  -- switch to indexes with scale & combine last + first
		AND (last_name ILIKE $2 OR $2 = '')
		AND (position = $3 OR $3 = '')
		AND (current_team_id = $4 OR $4 = 0)
		ORDER BY %s %s, id ASC
		LIMIT $5 OFFSET $6`, filters.sortColumn(), filters.SortDirection())

	args := []any{
		pq.FirstName,
		pq.LastName,
		pq.Position,
		pq.CurrentTeamID,
		filters.limit(),
		filters.offset(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}
	defer rows.Close()

	totalRecords := 0
	players := []*PlayerWithTeam{}

	for rows.Next() {
		var p PlayerWithTeam
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
			&p.TeamFullName,
			&p.TeamShortName,
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
