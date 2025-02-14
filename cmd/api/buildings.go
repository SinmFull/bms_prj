package main

import (
	"net/http"

	"github.com/SinmFull/BMS_prj/internal/data"
	"github.com/SinmFull/BMS_prj/internal/validator"
)

func (app *application) addNewBuilding(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Location string `json:"location"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	v := validator.New()
	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) <= 100, "name", "must not be more than 100 bytes long")
	v.Check(input.Location != "", "location", "must be provided")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	admin := app.contextGetUser(r)

	//Get the group id by email of admin
	ug, _ := app.models.UserGroups.Get(admin.Email)

	newBuilding := data.Building{GroupID: ug.ID, Name: input.Name, Location: input.Location}
	app.models.Buildings.Insert(&newBuilding)

	err = app.writeJSON(w, http.StatusCreated, envelope{"added_building": newBuilding}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showBuildings(w http.ResponseWriter, r *http.Request) {
	user := app.contextGetUser(r)
	if user.IsAnonymous() {
		app.invalidAuthenticationTokenResponse(w, r)
		return
	}

	buildings, err := app.models.Buildings.Get(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"buildings": buildings}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
