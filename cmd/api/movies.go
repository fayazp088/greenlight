package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fayazp088/greenlight/internal/models"
	"github.com/gin-gonic/gin"
)

func (app *application) createMovieHandler(c *gin.Context) {
	fmt.Fprintln(c.Writer, "create a new movie")
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
