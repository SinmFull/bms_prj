package main

import (
	"errors"
	"net/http"

	"github.com/SinmFull/BMS_prj/internal/data"
)

func (a *application) addSensorType(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
		Unit string `json:"unit"`
	}
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}

	sensorType := &data.SensorType{
		Name: input.Name,
		Unit: input.Unit,
	}
	err = a.models.SensorTypes.Insert(sensorType)
	if err != nil {
		if errors.Is(err, data.ErrSensorTypeExist) {
			a.badRequestResponse(w, r, err)
			return
		}
		a.serverErrorResponse(w, r, err)
		return
	}
	a.writeJSON(w, http.StatusCreated, envelope{"sensor_type": sensorType}, nil)
}

func (a *application) getSensorTypes(w http.ResponseWriter, r *http.Request) {
	user := a.contextGetUser(r)
	if user.IsAnonymous() {
		a.invalidAuthenticationTokenResponse(w, r)
		return
	}
	sensorTypes, err := a.models.SensorTypes.GetAll()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	a.writeJSON(w, http.StatusOK, envelope{"sensor_types": sensorTypes}, nil)
}