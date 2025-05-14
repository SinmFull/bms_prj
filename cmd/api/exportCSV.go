package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func (a *application) exportCSVHandler(w http.ResponseWriter, r *http.Request) {
	user := a.contextGetUser(r)
	if user.IsAnonymous() {
		a.invalidAuthenticationTokenResponse(w, r)
		return
	}

	startTime := r.URL.Query().Get("start")
	endTime := r.URL.Query().Get("end")
	device_id := r.URL.Query().Get("device_id")
	var input struct {
		Device_id int    `json:"device_id"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	input.StartTime = startTime
	input.EndTime = endTime
	input.Device_id, _ = strconv.Atoi(device_id)
	// err := a.readJSON(w, r, &input)
	// if err != nil {
	// 	a.badRequestResponse(w, r, err)
	// 	return
	// }

	if input.StartTime == "" || input.EndTime == "" {
		a.badRequestResponse(w, r, errors.New("start_time and end_time are required"))
		return
	}

	// Parse the time
	start, err := time.Parse("2006-01-02 15:04:05", input.StartTime)
	if err != nil {
		a.badRequestResponse(w, r, errors.New("start_time format error"))
		return
	}
	end, err := time.Parse("2006-01-02 15:04:05", input.EndTime)
	if err != nil {
		a.badRequestResponse(w, r, errors.New("end_time format error"))
		return
	}

	data, _ := a.models.SensorValue.GetBetweenTime(start, end, input.Device_id)

	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "attachment; filename=sensor_data.csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	writer.Write([]string{"ID", "Value", "Timestamp"})

	for _, row := range data {
		writer.Write([]string{
			fmt.Sprintf("%d", row.ID),
			row.Value,
			row.RecordedAt.Format("2006-01-02 15:04:05"),
		})
	}

}
