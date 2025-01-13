package main

import (
	"net/http"
	"time"

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

	app.writeJSON(c, http.StatusOK, envelope{"create_movies": input}, nil)
}

func (app *application) showMovieHandler(c *gin.Context) {
	id, err := app.readIDParam(c)

	if err != nil || id < 1 {
		app.notFoundResponse(c)
		return
	}

	movie := models.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "horror"},
		Version:   1,
	}

	app.writeJSON(c, http.StatusOK, envelope{"movie": movie}, nil)
}
