package main

import (
	"fmt"

	"github.com/SinmFull/BMS_prj/internal/data"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func (a *application) mqttMessageHandler(client mqtt.Client, msg mqtt.Message) {

	var sensorValue data.SensorValue
	sensorValue.SensorDeviceID = 1
	sensorValue.Value = string(msg.Payload())
	// json.Unmarshal(msg.Payload(), &sensorValue)
	// fmt.Println(sensorValue)
	// fmt.Println("Message Received: ", msg.Payload())
	err := a.models.SensorValue.Insert(&sensorValue)
	if err != nil {
		fmt.Println(err)
	}
}
