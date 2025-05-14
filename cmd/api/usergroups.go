package main

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/SinmFull/BMS_prj/internal/data"
)

func (app *application) addMemberToGroupHandler(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}
	var input struct {
		Email string `json:"email"`
	}
	app.readJSON(w, r, &input)
	var ug *data.UserGroup
	ug, _ = app.models.UserGroups.Get(user.Email)

	member, err := app.models.Users.GetByEmail(input.Email)
	if errors.Is(err, sql.ErrNoRows) || member == nil {
		app.userNotFoundResponse(w, r)
		return
	}
	ugm := &data.UserGroupMembers{
		UserID:  member.ID,
		GroupId: ug.ID,
		Role:    "member",
	}
	err = app.models.UserGroupMembers.Insert(ugm)
	if err != nil {
		if errors.Is(err, data.ErrAlreadyExists) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "member added to group"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (a *application) getGroupMembersHandler(w http.ResponseWriter, r *http.Request) {
	user := a.contextGetUser(r)
	if user.IsAnonymous() {
		a.invalidAuthenticationTokenResponse(w, r)
		return
	}
	if user.Role != "Admin" {
		a.notPermittedResponse(w, r)
		return
	}
	ug, err := a.models.UserGroups.Get(user.Email)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}

	users, err := a.models.UserGroupMembers.GetMembers(*ug)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	a.writeJSON(w, http.StatusOK, envelope{"members": users}, nil)
}