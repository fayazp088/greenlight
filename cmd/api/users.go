package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/fayazp088/greenlight/internal/models"
	"github.com/fayazp088/greenlight/internal/validator"
	"github.com/gin-gonic/gin"
)

func (app *application) registerUserHandler(c *gin.Context) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(c, &input)
	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	user := &models.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)

	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	v := validator.New()

	if models.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	err = app.models.User.Insert(user)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(c, v.Errors)
		default:
			app.serverErrorResponse(c, err)
		}

		return
	}

	// After the user record has been created in the database, generate a new activation
	// token for the user.
	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, models.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	app.background(func() {
		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
			return
		}
	})

	app.writeJSON(c, http.StatusCreated, envelope{"user": user}, nil)
}

func (app *application) activateUserHandler(c *gin.Context) {
	var input struct {
		TokenPlaintext string `json:"token"`
	}

	err := app.readJSON(c, &input)

	if err != nil {
		app.badRequestResponse(c, err)
		return
	}

	v := validator.New()
	if models.ValidateTokenPlaintext(v, input.TokenPlaintext); !v.Valid() {
		app.failedValidationResponse(c, v.Errors)
		return
	}

	user, err := app.models.User.GetForToken(models.ScopeActivation, input.TokenPlaintext)

	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(c, v.Errors)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	user.Activated = true

	err = app.models.User.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrEditConflict):
			app.editConflictResponse(c)
		default:
			app.serverErrorResponse(c, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(models.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(c, err)
		return
	}

	app.writeJSON(c, http.StatusOK, envelope{"user": user}, nil)
}
