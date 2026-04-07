package stats

type GoalieStatSet struct {
	CurrentSeason GoalieStats `json:"current_season,omitzero"`
	CareerTotals  GoalieStats `json:"career_totals,omitzero"`
}
