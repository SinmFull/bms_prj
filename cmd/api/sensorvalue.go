package main

import "net/http"

func (a *application) getNewestSensorReading(w http.ResponseWriter, r *http.Request) {
	user := a.contextGetUser(r)
	if user.IsAnonymous() {
		a.invalidAuthenticationTokenResponse(w, r)
		return
	}

	var input struct {
		SensorDeviceID int64 `json:"sensor_device_id"`
	}
	err := a.readJSON(w, r, &input)
	if err != nil {
		a.badRequestResponse(w, r, err)
		return
	}
	sensorValues, err := a.models.SensorValue.GetNowForDevice(input.SensorDeviceID)
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	a.writeJSON(w, http.StatusOK, envelope{"sensor_values": sensorValues}, nil)
}

func (a *application) getNewestForAllSensors(w http.ResponseWriter, r *http.Request) {
	sensorValues, err := a.models.SensorValue.GetNowForAllDevices()
	if err != nil {
		a.serverErrorResponse(w, r, err)
		return
	}
	a.writeJSON(w, http.StatusOK, envelope{"sensor_values": sensorValues}, nil)

}
