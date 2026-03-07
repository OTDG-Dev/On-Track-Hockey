package main

import (
	"errors"
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

// createEventParticipantHandler handles the creation of a new event participant
func (app *application) createEventParticipantHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		GameEventID int    `json:"game_event_id"`
		PlayerID    int    `json:"player_id"`
		Role        string `json:"role"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate required fields
	if input.GameEventID <= 0 {
		app.badRequestResponse(w, r, errors.New("game_event_id is required and must be positive"))
		return
	}
	if input.PlayerID <= 0 {
		app.badRequestResponse(w, r, errors.New("player_id is required and must be positive"))
		return
	}
	if input.Role == "" {
		app.badRequestResponse(w, r, errors.New("role is required"))
		return
	}

	participant := data.EventParticipant{
		GameEventID: input.GameEventID,
		PlayerID:    input.PlayerID,
		Role:        input.Role,
	}

	err = app.models.EventParticipants.Insert(&participant)
	if err != nil {
		app.handleDatabaseError(w, r, err)
		return
	}

	app.writeJSON(w, http.StatusCreated, envelope{"event_participant": participant}, nil)
}

// showEventParticipantHandler handles retrieving a single event participant by ID
func (app *application) showEventParticipantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	participant, err := app.models.EventParticipants.GetWithDetails(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.handleDatabaseError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"event_participant": participant}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listEventParticipantsByGameEventHandler handles retrieving all participants for a specific game event
func (app *application) listEventParticipantsByGameEventHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	participants, err := app.models.EventParticipants.ListByGameEvent(id)
	if err != nil {
		app.handleDatabaseError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"event_participants": participants}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// listEventParticipantsByPlayerHandler handles retrieving all event participants for a specific player
func (app *application) listEventParticipantsByPlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	participants, err := app.models.EventParticipants.ListByPlayer(id)
	if err != nil {
		app.handleDatabaseError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"event_participants": participants}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateEventParticipantHandler handles updating an existing event participant
func (app *application) updateEventParticipantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	participant, err := app.models.EventParticipants.Get(id)
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
		PlayerID *int    `json:"player_id"`
		Role     *string `json:"role"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.PlayerID != nil {
		participant.PlayerID = *input.PlayerID
	}
	if input.Role != nil {
		participant.Role = *input.Role
	}

	err = app.models.EventParticipants.Update(participant)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.handleDatabaseError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"event_participant": participant}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteEventParticipantHandler handles deleting an event participant
func (app *application) deleteEventParticipantHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.EventParticipants.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.handleDatabaseError(w, r, err)
		}
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"message": "event participant successfully deleted"}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
