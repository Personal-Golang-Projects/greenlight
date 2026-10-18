package main

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"greenlight.chitranshoo.net/internal/models"
	"greenlight.chitranshoo.net/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	// separate input struct so clients can't set ID, CreatedAt or Version
	var input struct {
		Title  string   `json:"title"`
		Year   int      `json:"year"`
		Genres []string `json:"genres"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		var badRequestErr *badRequestError
		switch {
		case errors.As(err, &badRequestErr):
			app.badRequestResponse(w, r, err)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	movie := &models.Movie{
		Title:  input.Title,
		Year:   input.Year,
		Genres: input.Genres,
	}

	v := validator.New()
	if models.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.notFoundResponse(w, r)
		return
	}

	movie := models.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Avatar",
		Genres:    []string{"science-fiction"},
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
