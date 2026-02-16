package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

func (app *application) showTeamHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	team, err := app.models.Teams.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"team": team}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createTeamHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName   string `json:"full_name"`
		ShortName  string `json:"short_name"`
		DivisionID int    `json:"division_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	team := &data.Team{
		FullName:   input.FullName,
		ShortName:  input.ShortName,
		DivisionID: input.DivisionID,
	}

	v := validator.New()

	if data.ValidateTeam(v, team); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Teams.Insert(team)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/players/%d", team.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"team": team}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Teams.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "team successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := app.models.Teams.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"teams": teams}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
