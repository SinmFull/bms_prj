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
	if errors.Is(err, sql.ErrNoRows) {
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
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "member added to group"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
