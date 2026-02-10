package data

import (
	"encoding/json"
	"errors"
)

type Player struct {
	PlayerID      int  `json:"player_id"`
	IsActive      bool `json:"is_active"`
	CurrentTeamId int  `json:"current_team_id"`

	FirstName     string         `json:"first_name"`
	LastName      string         `json:"last_name"`
	SweaterNumber uint8          `json:"sweater_number"`
	Position      Position       `json:"position"`
	BirthDate     BirthDate      `json:"birth_date"`
	BirthCountry  string         `json:"birth_country"`
	Headshot      string         `json:"headshot,omitzero"`
	ShootsCatches ShootsCatches  `json:"shoots_catches,omitzero"`
	SkaterStats   *SkaterStatSet `json:"skater_stats,omitzero"`
	GoalieStats   *GoalieStatSet `json:"goalie_stats,omitzero"`
}

type SeasonSplit[T any] struct {
	RegularSeason T `json:"regular_season,omitzero"`
	Playoffs      T `json:"playoffs,omitzero"`
}

type SkaterStatSet struct {
	CurrentStats SeasonSplit[SkaterStats] `json:"current_stats,omitzero"`
	CareerTotals SeasonSplit[SkaterStats] `json:"career_totals,omitzero"`
}

type GoalieStatSet struct {
	CurrentStats SeasonSplit[GoalieStats] `json:"current_stats,omitzero"`
	CareerTotals SeasonSplit[GoalieStats] `json:"career_totals,omitzero"`
}

var ErrInvalidPositionFormat = errors.New("invalid position format, expected: C|LW|RW|D|G")

type Position string

const (
	PositionC  Position = "C"
	PositionLW Position = "LW"
	PositionRW Position = "RW"
	PositionD  Position = "D"
	PositionG  Position = "G"
)

func (p *Position) UnmarshalJSON(jsonValue []byte) error {
	var s string
	if err := json.Unmarshal(jsonValue, &s); err != nil {
		return ErrInvalidPositionFormat
	}

	pos := Position(s)
	switch pos {
	case PositionC, PositionLW, PositionRW, PositionD, PositionG:
		*p = pos
		return nil
	}

	return ErrInvalidPositionFormat
}

var ErrInvalidShootCatches = errors.New("invalid position format, expected: L|R")

type ShootsCatches string

const (
	ShootsCatchesL ShootsCatches = "L"
	ShootsCatchesR ShootsCatches = "R"
)

func (sc *ShootsCatches) UnmarshalJSON(jsonValue []byte) error {
	var s string
	if err := json.Unmarshal(jsonValue, &s); err != nil {
		return ErrInvalidPositionFormat
	}

	shct := ShootsCatches(s)
	switch shct {
	case ShootsCatchesL, ShootsCatchesR:
		*sc = shct
		return nil
	}

	return ErrInvalidShootCatches
}
