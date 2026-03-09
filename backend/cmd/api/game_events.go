package main

import (
	"errors"
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

func (app *application) createGameEventHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GameID    int    `json:"game_id"`
		Period    int    `json:"period"`
		Clock     int    `json:"clock"`
		EventType string `json:"event_type"`
		TeamID    int    `json:"team_id"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	event := data.GameEvent{
		GameID:    input.GameID,
		Period:    input.Period,
		Clock:     input.Clock,
		EventType: input.EventType,
		TeamID:    input.TeamID,
	}

	v := validator.New()
	if data.ValidateGameEvent(v, &event); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.models.GameEvents.Insert(event); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"game_events": event}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showGameEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	event, err := app.models.GameEvents.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"game_events": event}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
