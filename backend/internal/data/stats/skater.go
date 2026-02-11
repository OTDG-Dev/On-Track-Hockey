package stats

type SkaterStatSet struct {
	CurrentStats SeasonSplit[SkaterStats] `json:"current_stats,omitzero"`
	CareerTotals SeasonSplit[SkaterStats] `json:"career_totals,omitzero"`
}
