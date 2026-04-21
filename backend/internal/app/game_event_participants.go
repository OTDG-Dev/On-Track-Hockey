package app

import (
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"
)

func (app *Application) createGameEventParticipantHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Role     string `json:"role"`
		PlayerID int    `json:"player_id"`
	}

	eventID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	part := data.GameEventParticipant{
		Role:     input.Role,
		PlayerID: input.PlayerID,
		EventID:  eventID,
	}

	v := validator.New()

	if data.ValidateGameEventParticiant(v, &part); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if err := app.Models.GameEventParticipants.Insert(part); err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusCreated, envelope{"game_event_participant": part}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) listGameEventParticipantsHandler(w http.ResponseWriter, r *http.Request) {
	eventID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	participants, err := app.Models.GameEventParticipants.GetByEvent(eventID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, envelope{"game_event_participants": participants}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
