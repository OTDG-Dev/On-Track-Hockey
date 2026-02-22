package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

func (app *application) showPlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	player, err := app.models.Players.GetWithTeam(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listPlayersHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName     string
		LastName      string
		Position      string
		CurrentTeamId int
		// WIP build this out & combine first/lastname search
		data.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.FirstName = app.readString(qs, "first_name", "")
	input.LastName = app.readString(qs, "last_name", "")
	input.Position = strings.ToUpper(app.readString(qs, "position", ""))

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id") // fallback to id

	input.Filters.SortSafeList = []string{
		"id", "-id",
		"first_name", "-first_name",
		"last_name", "-last_name",
		"position", "-position",
	}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	players, metadata, err := app.models.Players.GetAllWithTeam(input.FirstName, input.LastName, input.Position, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"players": players, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPlayerHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IsActive      bool `json:"is_active"`
		CurrentTeamId int  `json:"current_team_id"`

		FirstName     string             `json:"first_name"`
		LastName      string             `json:"last_name"`
		SweaterNumber uint8              `json:"sweater_number"`
		Position      data.Position      `json:"position"`
		BirthDate     data.BirthDate     `json:"birth_date"`
		BirthCountry  string             `json:"birth_country"`
		Headshot      string             `json:"headshot,omitzero"`
		ShootsCatches data.ShootsCatches `json:"shoots_catches,omitzero"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	player := &data.Player{
		IsActive:      input.IsActive,
		CurrentTeamID: input.CurrentTeamId,
		FirstName:     input.FirstName,
		LastName:      input.LastName,
		SweaterNumber: input.SweaterNumber,
		Position:      input.Position,
		BirthDate:     input.BirthDate,
		BirthCountry:  input.BirthCountry,
		Headshot:      input.Headshot,
		ShootsCatches: input.ShootsCatches,
	}

	v := validator.New()

	if data.ValidatePlayer(v, player); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Players.Insert(player)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("current_team_id", "team not found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/players/%d", player.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"player": player}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	player, err := app.models.Players.Get(id)
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
		IsActive      *bool `json:"is_active"`
		CurrentTeamId *int  `json:"current_team_id"`

		FirstName     *string             `json:"first_name"`
		LastName      *string             `json:"last_name"`
		SweaterNumber *uint8              `json:"sweater_number"`
		Position      *data.Position      `json:"position"`
		BirthDate     *data.BirthDate     `json:"birth_date"`
		BirthCountry  *string             `json:"birth_country"`
		Headshot      *string             `json:"headshot,omitzero"`
		ShootsCatches *data.ShootsCatches `json:"shoots_catches,omitzero"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.IsActive != nil {
		player.IsActive = *input.IsActive
	}
	if input.CurrentTeamId != nil {
		player.CurrentTeamID = *input.CurrentTeamId
	}
	if input.FirstName != nil {
		player.FirstName = *input.FirstName
	}
	if input.LastName != nil {
		player.LastName = *input.LastName
	}
	if input.SweaterNumber != nil {
		player.SweaterNumber = *input.SweaterNumber
	}
	if input.Position != nil {
		player.Position = *input.Position
	}
	if input.BirthDate != nil {
		player.BirthDate = *input.BirthDate
	}
	if input.BirthCountry != nil {
		player.BirthCountry = *input.BirthCountry
	}
	if input.Headshot != nil {
		player.Headshot = *input.Headshot
	}
	if input.ShootsCatches != nil {
		player.ShootsCatches = *input.ShootsCatches
	}

	v := validator.New()

	data.ValidatePlayer(v, player)

	err = app.models.Players.Update(player)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"player": player}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Players.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "player successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
