package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	resp := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}

	err := app.writeJSON(w, http.StatusOK, resp, nil)
	if err != nil {
		app.logger.Error(err.Error())
		http.Error(w, "The server encountered a problem", http.StatusInternalServerError)
	}
}
