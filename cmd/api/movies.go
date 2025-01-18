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
	Title   *string       `json:"title"`
	Year    *int32        `json:"year"`
	Runtime *data.Runtime `json:"runtime"`
	Genres  []string      `json:"genres"`
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

	if updateMovie.Title != nil {
		movie.Title = *updateMovie.Title
	}

	if updateMovie.Genres != nil {
		movie.Genres = updateMovie.Genres
	}

	if updateMovie.Runtime != nil {
		movie.Runtime = *updateMovie.Runtime
	}

	if updateMovie.Year != nil {
		movie.Year = *updateMovie.Year
	}

	v := validator.New()

	if models.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.Movies.Update(movie)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.serverErrorResponse(c, err)
			return
		}
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

func (app *application) listMoviesHandler(c *gin.Context) {
	var input struct {
		Title  string   `form:"title"`
		Genres []string `form:"genres"`
		data.Filters
	}

	if err := c.BindQuery(&input); err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	if input.Page == 0 {
		input.Page = 1
	}
	if input.PageSize == 0 {
		input.PageSize = 20
	}

	if input.Sort == "" {
		input.Sort = "id"
	}

	v := validator.New()

	input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	if input.Genres == nil {
		input.Genres = []string{}
	}

	movies, metaData, err := app.models.Movies.List(input.Title, input.Genres, input.Filters)

	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"movies": movies, "meta_data": metaData}, nil)
}
