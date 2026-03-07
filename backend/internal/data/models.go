package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Players          PlayerModel
	Teams            TeamModel
	Division         DivisonModel
	League           LeagueModel
	Games            GameModel
	GameEvents       GameEventModel
	EventParticipants EventParticipantModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Players:           PlayerModel{DB: db},
		Teams:             TeamModel{DB: db},
		Division:          DivisonModel{DB: db},
		League:            LeagueModel{DB: db},
		Games:             GameModel{DB: db},
		GameEvents:        GameEventModel{DB: db},
		EventParticipants: EventParticipantModel{DB: db},
	}
}
