package stats

type SeasonSplit[T any] struct {
	RegularSeason T `json:"regular_season,omitzero"`
	Playoffs      T `json:"playoffs,omitzero"`
}
