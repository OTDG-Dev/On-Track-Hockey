package stats

type SkaterStatSet struct {
	CurrentSeason SkaterStats `json:"current_season"`
	CareerTotals  SkaterStats `json:"career_totals"`
}
