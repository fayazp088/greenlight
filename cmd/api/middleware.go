package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func (app *application) recoverPanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Connection", "close")
				app.serverErrorResponse(c, fmt.Errorf("%v", err))
				c.Abort()
			}
		}()
		c.Next()
	}
}

type inputValidationErrors struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (app *application) inputValidation() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Proceed with the request first
		c.Next()

		// Check for any validation errors after the handler has run
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				// Check if it's a validation error
				if validationErr, ok := err.Err.(validator.ValidationErrors); ok {
					// Format the validation error into a structured array response
					var errorMessages []inputValidationErrors
					for _, fieldErr := range validationErr {
						field := fieldErr.Field() // Extract the field name
						message := ""

						// Map validation tags to user-friendly messages
						switch fieldErr.Tag() {
						case "required":
							message = field + " is required."
						case "min":
							message = field + " must be at least " + fieldErr.Param() + "."
						case "max":
							message = field + " cannot exceed " + fieldErr.Param() + "."
						case "email":
							message = field + " must be a valid email address."
						// Add more cases as needed for other tags
						default:
							message = fieldErr.Error()
						}

						// Append the error object to the array
						errorMessages = append(errorMessages, inputValidationErrors{
							Key:   field,
							Value: message,
						})
					}

					app.writeJSON(c, http.StatusBadRequest, envelope{
						"status":  "error",
						"message": "Validation failed",
						"errors":  errorMessages,
					}, nil)

					return
				}
			}
		}
	}
}
