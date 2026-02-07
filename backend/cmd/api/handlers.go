package main

import (
	"net/http"
	"time"

	"github.com/OTDG-Dev/On-Track-Hockey/backend/internal/data"
)

func (app *application) showPlayerHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	// tmp dummy data
	if id != 1 {
		w.Write([]byte(`{"status": "error"}`))
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
		BirthDate: data.Date{
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
		app.logger.Error(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
