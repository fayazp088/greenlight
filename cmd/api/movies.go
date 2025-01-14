package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/fayazp088/greenlight/internal/data"
	"github.com/fayazp088/greenlight/internal/models"
	"github.com/fayazp088/greenlight/internal/validator"
	"github.com/gin-gonic/gin"
)

type CreateMovieInput struct {
	Title   string       `json:"title"`
	Year    int32        `json:"year"`
	Runtime data.Runtime `json:"runtime"`
	Genres  []string     `json:"genres"`
}

type UpdateMovieInput struct {
	Title   string       `json:"title"`
	Year    int32        `json:"year"`
	Runtime data.Runtime `json:"runtime"`
	Genres  []string     `json:"genres"`
}

func (app *application) createMovieHandler(c *gin.Context) {

	var input CreateMovieInput
	err := app.readJSON(c, &input)

	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	movie := &models.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if models.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.Movies.Insert(movie)

	if err != nil {
		app.errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	c.Header("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	app.writeJSON(c, http.StatusOK, envelope{"create_movies": movie}, nil)
}

func (app *application) updateMovieHandler(c *gin.Context) {
	id, err := app.readIDParam(c)

	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	if id < 0 {
		app.errorResponse(c, http.StatusBadRequest, models.ErrRecordNotFound)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	var updateMovie UpdateMovieInput

	err = app.readJSON(c, &updateMovie)

	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	movie.Title = updateMovie.Title
	movie.Genres = updateMovie.Genres
	movie.Runtime = updateMovie.Runtime
	movie.Year = updateMovie.Year

	err = app.models.Movies.Update(movie)

	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"movie": movie}, nil)

}

func (app *application) showMovieHandler(c *gin.Context) {
	id, err := app.readIDParam(c)

	if err != nil || id < 1 {
		app.notFoundResponse(c)
		return
	}

	movie, err := app.models.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"movie": movie}, nil)
}

func (app *application) deleteMovieHandler(c *gin.Context) {
	id, err := app.readIDParam(c)

	if err != nil || id < 1 {
		app.notFoundResponse(c)
		return
	}

	err = app.models.Movies.Delete(id)

	if err != nil {
		app.errorResponse(c, http.StatusInternalServerError, err)
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"message": "movie successfully deleted"}, nil)
}
