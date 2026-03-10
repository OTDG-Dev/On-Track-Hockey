package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

func (app *application) createGameHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		HomeTeamId int       `json:"home_team_id"`
		AwayTeamID int       `json:"away_team_id"`
		StartTime  time.Time `json:"start_time"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	game := data.Game{
		HomeTeamID: input.HomeTeamId,
		AwayTeamID: input.AwayTeamID,
		StartTime:  input.StartTime,
	}

	err = app.models.Games.Insert(&game)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"game": game}, nil)
}

func (app *application) showGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	game, err := app.models.Games.GetView(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
