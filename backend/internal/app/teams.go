package app

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/validator"
)

func (app *Application) showTeamHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	team, err := app.Models.Teams.Get(id)
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

func (app *Application) createTeamHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FullName   string `json:"full_name"`
		ShortName  string `json:"short_name"`
		DivisionID int    `json:"division_id"`
		IsActive   bool   `json:"is_active"`
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
		IsActive:   input.IsActive,
	}

	v := validator.New()

	if data.ValidateTeam(v, team); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.Models.Teams.Insert(team)

	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("division_id", "division not found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/teams/%d", team.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"team": team}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *Application) updateTeamHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	team, err := app.Models.Teams.Get(id)
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
		FullName   *string `json:"full_name"`
		ShortName  *string `json:"short_name"`
		DivisionID *int    `json:"division_id"`
		IsActive   *bool   `json:"is_active"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.FullName != nil {
		team.FullName = *input.FullName
	}
	if input.ShortName != nil {
		team.ShortName = *input.ShortName
	}
	if input.DivisionID != nil {
		team.DivisionID = *input.DivisionID
	}
	if input.IsActive != nil {
		team.IsActive = *input.IsActive
	}

	err = app.Models.Teams.Update(team)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
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

func (app *Application) deleteTeamHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.Models.Teams.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrFKeyViolation):
			app.badRequestResponse(w, r, err)
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

func (app *Application) listTeamsHandler(w http.ResponseWriter, r *http.Request) {
	teams, err := app.Models.Teams.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"teams": teams}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
