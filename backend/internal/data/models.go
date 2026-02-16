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
	Players PlayerModel
	Teams   TeamModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Players: PlayerModel{DB: db},
		Teams:   TeamModel{DB: db},
	}
}
