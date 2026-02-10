package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data/validator"
)

func (app *application) showPlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// tmp dummy data
	if id != 1 {
		app.notFoundResponse(w, r)
		return
	}

	skater := data.Player{
		PlayerID:      1,
		IsActive:      false,
		CurrentTeamId: 456,
		FirstName:     "Connor",
		LastName:      "McDavid",
		SweaterNumber: 97,
		Position:      data.PositionC,
		BirthDate: data.BirthDate{
			Time: time.Date(1997, 1, 13, 0, 0, 0, 0, time.UTC),
		},
		Headshot: "https://assets/path/to/headshot.png",
		SkaterStats: &data.SkaterStatSet{
			CurrentStats: data.SeasonSplit[data.SkaterStats]{
				RegularSeason: data.SkaterStats{
					BasicStats: data.BasicStats{
						GamesPlayed: 32,
						Assists:     14,
						Goals:       12,
						PIM:         6,
						TOI: data.TimeOnIce{
							Duration: 23*time.Minute + 14*time.Second,
						},
					},
					Shots:             24,
					PlusMinus:         4,
					OTGoals:           0,
					GameWinningGoals:  1,
					ShortHandedGoals:  0,
					ShortHandedPoints: 0,
					PowerPlayGoals:    1,
					PowerPlayPoints:   1,
				},
				Playoffs: data.SkaterStats{},
			},
			CareerTotals: data.SeasonSplit[data.SkaterStats]{
				RegularSeason: data.SkaterStats{
					BasicStats: data.BasicStats{
						GamesPlayed: 32,
						Assists:     14,
						Goals:       12,
						PIM:         6,
						TOI: data.TimeOnIce{
							Duration: 23*time.Minute + 14*time.Second,
						},
					},
					Shots:             24,
					PlusMinus:         4,
					OTGoals:           0,
					GameWinningGoals:  1,
					ShortHandedGoals:  0,
					ShortHandedPoints: 0,
					PowerPlayGoals:    1,
					PowerPlayPoints:   1,
				},
				Playoffs: data.SkaterStats{},
			},
		},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"player": skater}, nil)
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
		CurrentTeamId: input.CurrentTeamId,
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

	fmt.Fprintln(w, "Success!")
	fmt.Fprintf(w, "%v\n", input)
}
