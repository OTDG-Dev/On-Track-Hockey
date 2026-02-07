package data

type BasicStats struct {
	GamesPlayed int       `json:"games_played"`
	Assists     int       `json:"assists"`
	Goals       int       `json:"goals"`
	PIM         int       `json:"pim"`
	TOI         TimeOnIce `json:"time_on_ice"`
}

type SkaterStats struct {
	BasicStats `json:"basic_stats"`

	Shots int `json:"shots"`

	PlusMinus int `json:"plus_minus"`

	OTGoals          int `json:"ot_goals"`
	GameWinningGoals int `json:"game_winning_goals"`

	ShortHandedGoals  int `json:"shorthanded_goals"`
	ShortHandedPoints int `json:"shorthanded_points"`
	PowerPlayGoals    int `json:"powerplay_goals"`
	PowerPlayPoints   int `json:"powerplay_points"`
}

type GoalieStats struct {
	BasicStats `json:"basic_stats"`

	GamesStarted int `json:"games_started"`
	GoalsAgainst int `json:"goals_against"`
	// calc: GoalsAgainstAvg

	Wins     int `json:"wins"`
	Losses   int `json:"losses"`
	OTLosses int `json:"ot_losses"`

	ShotsAgainst int `json:"shots_against"`
	Shutouts     int `json:"shutouts"`
}

func (s SkaterStats) Points() int {
	return s.Goals + s.Assists
}

func (s SkaterStats) ShootingPctg() float64 {
	if s.Shots == 0 {
		return 0
	}
	return float64(s.Goals) / float64(s.Shots)
}

func (s GoalieStats) SavePctg() float64 {
	saves := s.ShotsAgainst - s.GoalsAgainst
	if saves == 0 {
		return 0
	}
	return float64(saves) / float64(s.ShotsAgainst)
}

func (s GoalieStats) GoalsAgainstAvg() float64 {
	if s.GamesPlayed == 0 {
		return 0
	}
	return float64(s.GoalsAgainst) / float64(s.GamesPlayed)
}
