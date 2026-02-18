package main

import (
	"fmt"
	"net/http"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

func (app *application) createDivisionHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		LeagueID int    `json:"league_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	div := &data.Division{
		LeagueID: input.LeagueID,
		Name:     input.Name,
	}

	v := validator.New()

	if data.ValidateDivision(v, div); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Division.Insert(div)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/divisions/%d", div.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"division": div}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listDivisionsHandler(w http.ResponseWriter, r *http.Request) {
	divs, err := app.models.Division.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"divisions": divs}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
