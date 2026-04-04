package app

import "net/http"

func (app *Application) updateStatsHandler(w http.ResponseWriter, r *http.Request) {
	updatedPlayers, err := app.Models.Players.RebuildStats()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{
		"message":         "player stats updated",
		"updated_players": updatedPlayers,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
