package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// send JSON error response when error is raised for a route
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// healthcheck api
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// create movies
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)

	// show a movie
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	// wrap panic recover middleware and return router instance
	return app.recoverPanic(router)
}
