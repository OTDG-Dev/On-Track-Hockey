package app

import "net/http"

func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.Config.Env,
		},
	}

	err := app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.Logger.Error(err.Error())
		http.Error(w, "The server encountered a problem", http.StatusInternalServerError)
	}
}
