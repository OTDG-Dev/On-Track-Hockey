package main

import (
	"errors"
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

// createGameEventHandler handles the creation of a new game event
func (app *application) createGameEventHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GameID    int    `json:"game_id"`
		Period    int    `json:"period"`
		Clock     string `json:"clock"`
		EventType string `json:"event_type"`
		TeamID    int    `json:"team_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate required fields
	if input.GameID <= 0 {
		app.badRequestResponse(w, r, errors.New("game_id is required and must be positive"))
		return
	}
	if input.Period <= 0 {
		app.badRequestResponse(w, r, errors.New("period is required and must be positive"))
		return
	}
	if input.Clock == "" {
		app.badRequestResponse(w, r, errors.New("clock is required"))
		return
	}
	if input.EventType == "" {
		app.badRequestResponse(w, r, errors.New("event_type is required"))
		return
	}
	if input.TeamID <= 0 {
		app.badRequestResponse(w, r, errors.New("team_id is required and must be positive"))
		return
	}

	event := data.GameEvent{
		GameID:    input.GameID,
		Period:    input.Period,
		Clock:     input.Clock,
		EventType: input.EventType,
		TeamID:    input.TeamID,
	}

	err = app.models.GameEvents.Insert(&event)
	if err != nil {
		app.handleDatabaseError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"game_event": event}, nil)
}

// showGameEventHandler handles retrieving a single game event by ID
func (app *application) showGameEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	event, err := app.models.GameEvents.GetWithDetails(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.handleDatabaseError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"game_event": event}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listGameEventsByGameHandler handles retrieving all game events for a specific game
func (app *application) listGameEventsByGameHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	events, err := app.models.GameEvents.ListByGame(id)
	if err != nil {
		app.handleDatabaseError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"game_events": events}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateGameEventHandler handles updating an existing game event
func (app *application) updateGameEventHandler(w http.ResponseWriter, r *http.Request) {
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
			app.handleDatabaseError(w, r, err)
		}
		return
	}

	var input struct {
		Period    *int    `json:"period"`
		Clock     *string `json:"clock"`
		EventType *string `json:"event_type"`
		TeamID    *int    `json:"team_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Period != nil {
		event.Period = *input.Period
	}
	if input.Clock != nil {
		event.Clock = *input.Clock
	}
	if input.EventType != nil {
		event.EventType = *input.EventType
	}
	if input.TeamID != nil {
		event.TeamID = *input.TeamID
	}

	err = app.models.GameEvents.Update(event)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.handleDatabaseError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"game_event": event}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteGameEventHandler handles deleting a game event
func (app *application) deleteGameEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.GameEvents.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.handleDatabaseError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "game event successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
