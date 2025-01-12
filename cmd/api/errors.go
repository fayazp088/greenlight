package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// The logError() method is a generic helper for logging an error message along
// with the current request method and URL as attributes in the log entry.
func (app *application) logError(c *gin.Context, err error) {
	var (
		method = c.Request.Method
		uri    = c.Request.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

// The errorResponse() method is a generic helper for sending JSON-formatted error
// messages to the client with a given status code. Note that we're using the any
// type for the message parameter, rather than just a string type, as this gives us
// more flexibility over the values that we can include in the response.
func (app *application) errorResponse(c *gin.Context, status int, message any) {
	env := envelope{"error": message}
	// Write the response using the writeJSON() helper. If this happens to return an
	// error then log it, and fall back to sending the client an empty response with a
	// 500 Internal Server Error status code.
	app.writeJSON(c, status, env, nil)
}

func (app *application) serverErrorResponse(c *gin.Context, err error) {
	app.logError(c, err)
	message :=
		"the server encountered a problem and could not process your request"
	app.errorResponse(c, http.StatusInternalServerError, message)

}

func (app *application) notFoundResponse(c *gin.Context) {
	message := "the requested resource could not be found"
	app.errorResponse(c, http.StatusNotFound, message)
}

// status code and JSON response to the client.
func (app *application) methodNotAllowedResponse(c *gin.Context) {
	message := fmt.Sprintf("the %s method is not supported for this resource", c.Request.Method)
	app.errorResponse(c, http.StatusMethodNotAllowed, message)
}

func (app *application) badRequestResponse(c *gin.Context, err error) {
	app.errorResponse(c, http.StatusBadRequest, err.Error())
}
