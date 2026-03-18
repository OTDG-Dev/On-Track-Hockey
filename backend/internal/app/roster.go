package app

import "net/http"

func (app *Application) listRosterHandler(w http.ResponseWriter, r *http.Request) {
	teamID, err := app.readIDParam(r, "team_id")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	roster, err := app.Models.Roster.Get(teamID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusAccepted, envelope{"roster": roster}, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
