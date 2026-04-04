package app

import (
	"errors"
	"net/http"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

func (app *Application) createGameHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		HomeTeamId int       `json:"home_team_id"`
		AwayTeamID int       `json:"away_team_id"`
		StartTime  time.Time `json:"start_time"`
		IsFinished bool      `json:"is_finished"`
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
		IsFinished: input.IsFinished,
	}

	err = app.Models.Games.Insert(&game)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"game": game}, nil)
}

func (app *Application) showGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	game, err := app.Models.Games.GetView(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If a game is marked as finished, we need to recalculate player stats to reflect the final results of the game
	// This is a simple approach and should be optimizd to only recalculate stats for players involved in the game, but it works for now
	if game.IsFinished {
		_, err := app.Models.Players.RebuildStats()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) updateGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	game, err := app.Models.Games.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		HomeTeamID *int  `json:"home_team_id"`
		AwayTeamID *int  `json:"away_team_id"`
		IsFinished *bool `json:"is_finished"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.HomeTeamID != nil {
		game.HomeTeamID = *input.HomeTeamID
	}
	if input.AwayTeamID != nil {
		game.AwayTeamID = *input.AwayTeamID
	}
	if input.IsFinished != nil {
		game.IsFinished = *input.IsFinished
	}

	err = app.Models.Games.Update(game)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// If a game is marked as finished, we need to recalculate player stats to reflect the final results of the game
	// This is a simple approach and should be optimizd to only recalculate stats for players involved in the game, but it works for now
	if game.IsFinished {
		_, err := app.Models.Players.RebuildStats()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"game": game}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *Application) listGamesHandler(w http.ResponseWriter, r *http.Request) {
	games, err := app.Models.Games.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"games": games}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
