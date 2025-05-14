package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
)

var validSeconds = map[string]bool{
	"30":   true,
	"60":   true,
	"300":  true,
	"600":  true,
	"900":  true,
	"1200": true,
	"1800": true,
	"3600": true,
}

var validMinutes = map[string]bool{
	"1":    true,
	"5":    true,
	"10":   true,
	"15":   true,
	"20":   true,
	"30":   true,
	"60":   true,
	"1440": true,
}

type SetCommandMessage struct {
	OprID string `json:"oprid"`
	Cmd   string `json:"Cmd"`
	Value string `json:"value"`
	Types string `json:"types"`
}

func (a *application) generateRandomOprID() (string, error) {
	bytes := make([]byte, 16)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

func (a *application) setDataFreqHandler(w http.ResponseWriter, r *http.Request) {
	production_id := r.URL.Query().Get("production_id")
	time_level := r.URL.Query().Get("level")
	time_value := r.URL.Query().Get("value")

	if production_id == "" || time_level == "" || time_value == "" {
		a.badRequestResponse(w, r, errors.New("production_id, level and value are required"))
		return
	}

	if time_level == "0" {
		if _, ok := validSeconds[time_value]; !ok {
			a.badRequestResponse(w, r, errors.New("second value is not valid"))
			return
		}
	} else if time_level == "1" {
		if _, ok := validMinutes[time_value]; !ok {
			a.badRequestResponse(w, r, errors.New("minute value is not valid"))
			return
		}
	} else {
		a.badRequestResponse(w, r, errors.New("time level is not valid(0 for second, 1 for minute)"))
		return
	}

	oprID, _ := a.generateRandomOprID()

	mqtt_topic := "MQTT_COMMOD_SET_" + production_id[len(production_id)-8:]
	Cmd := "000" + time_level
	msg := SetCommandMessage{
		OprID: oprID,
		Cmd:   Cmd,
		Value: time_value,
		Types: "1",
	}

	payload, err := json.Marshal(msg)
	if err != nil {
		a.logger.PrintError(err, nil)
		return
	}

	token := a.mqttClient.Publish(mqtt_topic, 0, false, payload)
	token.Wait()

	a.writeJSON(w, http.StatusOK, envelope{"status": "ok"}, nil)
}
