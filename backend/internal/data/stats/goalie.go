package stats

type GoalieStatSet struct {
	CurrentStats SeasonSplit[GoalieStats] `json:"current_stats,omitzero"`
	CareerTotals SeasonSplit[GoalieStats] `json:"career_totals,omitzero"`
}
