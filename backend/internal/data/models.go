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
	Players  PlayerModel
	Teams    TeamModel
	Division DivisonModel
	League   LeagueModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Players:  PlayerModel{DB: db},
		Teams:    TeamModel{DB: db},
		Division: DivisonModel{DB: db},
		League:   LeagueModel{DB: db},
	}
}
