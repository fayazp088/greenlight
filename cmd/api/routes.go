package main

import "github.com/gin-gonic/gin"

func (app *application) routes() *gin.Engine {
	router := gin.Default()
	router.Use(app.recoverPanic())

	v1 := router.Group("/v1")
	{
		v1.GET("/health", app.Health)
		v1.POST("/movies", app.createMovieHandler)
		v1.GET("/movies/:id", app.showMovieHandler)
	}

	router.NoMethod(app.methodNotAllowedResponse)

	router.NoRoute(app.notFoundResponse)

	return router
}
