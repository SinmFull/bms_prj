package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/SinmFull/BMS_prj/internal/data"
	"github.com/SinmFull/BMS_prj/internal/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:  input.Name,
		Email: input.Email,
		Role:  data.UserRole("User"),
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) userLoginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()
	v.Check(input.Email != "", "email", "email must be provided")
	v.Check(input.Password != "", "password", "password must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, _ := app.models.Users.GetByEmail(input.Email)
	if user == nil {
		app.userNotFoundResponse(w, r)
		return
	}

	valid, _ := user.Password.Matches(input.Password)
	if !valid {
		app.invalidAccountResponse(w, r)
	}

	token, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.ScopeLogin)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"user": user, "authentication_token": token}, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
