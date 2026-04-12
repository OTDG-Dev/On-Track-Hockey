package stats

type GoalieStatSet struct {
	CurrentSeason GoalieStats `json:"current_season"`
	CareerTotals  GoalieStats `json:"career_totals"`
}
