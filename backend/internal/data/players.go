package data

type Position string

const (
	PositionC  Position = "C"
	PositionLW Position = "LW"
	PositionRW Position = "RW"
	PositionD  Position = "D"
	PositionG  Position = "G"
)

type ShootsCatches string

const (
	ShootsCatchesL ShootsCatches = "L"
	ShootsCatchesR ShootsCatches = "R"
)

type Player struct {
	PlayerID      int  `json:"player_id"`
	IsActive      bool `json:"is_active"`
	CurrentTeamId int  `json:"current_team_id"`

	FirstName     string        `json:"first_name"`
	LastName      string        `json:"last_name"`
	SweaterNumber uint8         `json:"sweater_number"`
	Position      Position      `json:"position"`
	BirthDate     Date          `json:"birth_date"`    // force format 1995-04-22
	BirthCountry  string        `json:"birth_country"` // "CAN"
	Headshot      string        `json:"headshot,omitzero"`
	ShootsCatches ShootsCatches `json:"shoots_catches,omitzero"`
	// "playerSlug": "philipp-grubauer-8475831", // optional, derive at call time
	SkaterStats *SkaterStatSet `json:"skater_stats,omitzero"`
	GoalieStats *GoalieStatSet `json:"goalie_stats,omitzero"`
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
