package main

import (
	"net/http"

	"github.com/SinmFull/BMS_prj/internal/data"
)

func (a *application) addOneSensorForBuilding(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BuildingID   int64  `json:"building_id"`
		SensorTypeID int64  `json:"sensor_type_id"`
		Name         string `json:"name"`
		Location     string `json:"location"`
	}
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	sensor := &data.SensorDevice{
		BuildingID:   input.BuildingID,
		SensorTypeID: input.SensorTypeID,
		Name:         input.Name,
		Location:     input.Location,
	}
	err = a.models.SensorDevices.Insert(sensor)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	a.writeJSON(w, http.StatusOK, envelope{"sensor": sensor}, nil)
}

func (a *application) getAllSensorsForOneBuilding(w http.ResponseWriter, r *http.Request) {
	var input struct {
		BuildingID int64 `json:"building_id"`
	}
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	sensors, err := a.models.SensorDevices.GetAllForBuilding(input.BuildingID)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	a.writeJSON(w, http.StatusOK, envelope{"sensors": sensors}, nil)
}
