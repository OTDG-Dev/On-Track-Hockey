package stats

type SkaterStatSet struct {
	CurrentSeason SkaterStats `json:"current_season,omitzero"`
	CareerTotals  SkaterStats `json:"career_totals,omitzero"`
}
