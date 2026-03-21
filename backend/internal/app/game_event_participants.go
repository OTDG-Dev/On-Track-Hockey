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

	eventID, err := app.readIDParam(r, "event_id")
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
