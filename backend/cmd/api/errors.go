package main

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/lib/pq"
)

// helper method for logging an error message along with current request
func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// JSON formatted error messages, using any type as the message
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{
		"error": message,
	}

	// Fall back to raw 500 if writing the error fails
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// Send a 500, log the request and send the user a generic JSON response
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// Send a 404, log the request and send the user a generic JSON response
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// Send a 405, log the request and send the user a generic JSON response
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the %s method is not supported for this resouce", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// Send a 400, log the request and send the user a generic JSON response
// If user input is bad
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusTooManyRequests, "rate limit exceeded")
}

// handleDatabaseError provides verbose error messages for database errors,
// particularly for foreign key constraint violations (Issue #73)
func (app *application) handleDatabaseError(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	// Check for PostgreSQL specific errors
	var pgErr *pq.Error
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23503": // foreign_key_violation
			message := app.formatForeignKeyError(pgErr)
			app.errorResponse(w, r, http.StatusConflict, message)
			return
		case "23505": // unique_violation
			message := fmt.Sprintf("A record with this %s already exists", pgErr.Constraint)
			app.errorResponse(w, r, http.StatusConflict, message)
			return
		case "23502": // not_null_violation
			message := fmt.Sprintf("Field %s cannot be empty", pgErr.Column)
			app.errorResponse(w, r, http.StatusBadRequest, message)
			return
		case "23514": // check_violation
			message := fmt.Sprintf("Value violates check constraint: %s", pgErr.Constraint)
			app.errorResponse(w, r, http.StatusBadRequest, message)
			return
		}
	}

	// Default to server error for unhandled database errors
	app.serverErrorResponse(w, r, err)
}

// formatForeignKeyError creates a human-readable error message for foreign key violations
func (app *application) formatForeignKeyError(pgErr *pq.Error) string {
	constraint := pgErr.Constraint
	detail := pgErr.Detail

	var message strings.Builder
	message.WriteString("Unable to complete operation: ")

	// Try to provide specific context based on the constraint
	switch {
	case strings.Contains(constraint, "game_id"):
		message.WriteString("the referenced game does not exist")
	case strings.Contains(constraint, "team_id"):
		message.WriteString("the referenced team does not exist")
	case strings.Contains(constraint, "player_id"):
		message.WriteString("the referenced player does not exist")
	case strings.Contains(constraint, "game_event_id"):
		message.WriteString("the referenced game event does not exist")
	case strings.Contains(constraint, "home_team_id"):
		message.WriteString("the referenced home team does not exist")
	case strings.Contains(constraint, "away_team_id"):
		message.WriteString("the referenced away team does not exist")
	case strings.Contains(constraint, "division_id"):
		message.WriteString("the referenced division does not exist")
	case strings.Contains(constraint, "league_id"):
		message.WriteString("the referenced league does not exist")
	default:
		message.WriteString("a related record does not exist")
	}

	// Add detail if available
	if detail != "" {
		message.WriteString(" (")
		message.WriteString(detail)
		message.WriteString(")")
	}

	return message.String()
}
