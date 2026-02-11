package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Players PlayerModel
}

func NewModel(db *sql.DB) Models {
	return Models{
		Players: PlayerModel{DB: db},
	}
}
