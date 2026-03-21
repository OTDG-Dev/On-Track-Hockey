package app

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/players", app.listPlayersHandler)
	router.HandlerFunc(http.MethodPost, "/v1/players", app.createPlayerHandler)
	router.HandlerFunc(http.MethodGet, "/v1/players/:id", app.showPlayerHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/players/:id", app.updatePlayerHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/players/:id", app.deletePlayerHandler)

	router.HandlerFunc(http.MethodGet, "/v1/teams", app.listTeamsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/teams", app.createTeamHandler)
	router.HandlerFunc(http.MethodGet, "/v1/teams/:id", app.showTeamHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/teams/:id", app.updateTeamHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/teams/:id", app.deleteTeamHandler)

	router.HandlerFunc(http.MethodGet, "/v1/roster/:team_id", app.listRosterHandler)

	router.HandlerFunc(http.MethodGet, "/v1/divisions", app.listDivisionsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/divisions", app.createDivisionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/divisions/:id", app.showDivisionHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/divisions/:id", app.updateDivisionHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/divisions/:id", app.deleteDivisionHandler)
	router.HandlerFunc(http.MethodGet, "/v1/divisions/:id/teams", app.listDivisionTeamHandler)

	router.HandlerFunc(http.MethodGet, "/v1/leagues", app.listLeaguesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/leagues", app.createLeagueHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues/:id", app.showLeagueHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/leagues/:id", app.updateLeagueHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/leagues/:id", app.deleteLeagueHandler)
	router.HandlerFunc(http.MethodGet, "/v1/leagues/:id/divisions", app.listLeagueDivisonHandler)

	router.HandlerFunc(http.MethodPost, "/v1/games", app.createGameHandler)
	router.HandlerFunc(http.MethodGet, "/v1/games/:id", app.showGameHandler)

	router.HandlerFunc(http.MethodPost, "/v1/games/:game_id/events", app.createGameEventHandler)
	router.HandlerFunc(http.MethodGet, "/v1/events/:id", app.showGameEventHandler)

	router.HandlerFunc(http.MethodPost, "/v1/events/:event_id/participants", app.createGameEventParticipantHandler)

	return app.recoverPanic(app.rateLimit(app.cors(router)))
}
