package data

import (
	"database/sql"
)

type Models struct {
	Players    PlayerModel
	Teams      TeamModel
	Division   DivisonModel
	League     LeagueModel
	Games      GameModel
	GameEvents GameEventModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Players:    PlayerModel{DB: db},
		Teams:      TeamModel{DB: db},
		Division:   DivisonModel{DB: db},
		League:     LeagueModel{DB: db},
		Games:      GameModel{DB: db},
		GameEvents: GameEventModel{DB: db},
	}
}
