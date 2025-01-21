package main

import "github.com/gin-gonic/gin"

func (app *application) routes() *gin.Engine {
	router := gin.Default()
	router.Use(app.inputValidation())
	router.Use(app.recoverPanic())
	router.Use(app.rateLimiter())

	v1 := router.Group("/v1")
	{
		v1.GET("/health", app.Health)
		v1.POST("/movies", app.createMovieHandler)
		v1.GET("/movies", app.listMoviesHandler)
		v1.PATCH("/movies/:id", app.updateMovieHandler)
		v1.GET("/movies/:id", app.showMovieHandler)
		v1.DELETE("/movies/:id", app.deleteMovieHandler)

		v1.POST("/users", app.registerUserHandler)
	}

	router.NoMethod(app.methodNotAllowedResponse)

	router.NoRoute(app.notFoundResponse)

	return router
}
