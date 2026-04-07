package stats

type SkaterStats struct {
	GamesPlayed int `json:"games_played"`
	Goals       int `json:"goals"`
	Assists     int `json:"assists"`
	Points      int `json:"points"`
	PIM         int `json:"pim"`
}

type GoalieStats struct {
	GamesPlayed  int `json:"games_played"`
	Wins         int `json:"wins"`
	Losses       int `json:"losses"`
	OTLosses     int `json:"ot_losses"`
	GoalsAgainst int `json:"goals_against"`
	Shutouts     int `json:"shutouts"`
}
