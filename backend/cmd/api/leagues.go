package main

import (
	"fmt"
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

func (app *application) createLeagueHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	league := &data.League{
		Name: input.Name,
	}

	v := validator.New()

	if data.ValidateLeague(v, league); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.League.Insert(league)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/leagues/%d", league.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"league": league}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listLeaguesHandler(w http.ResponseWriter, r *http.Request) {
	leagues, err := app.models.League.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"leagues": leagues}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
