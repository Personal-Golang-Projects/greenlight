package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// healthcheck api
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// create movies
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)

	// show a movie
	router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	// return router instance
	return router
}
