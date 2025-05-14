package main

import (
	"encoding/json"
	"strconv"
	"sync"

	"github.com/SinmFull/BMS_prj/internal/data"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTMessage struct {
	ID    string  `json:"id"`
	U0    float64 `json:"u0"`
	UP    float64 `json:"u+"`
	UN    float64 `json:"u-"`
	I0    float64 `json:"i0"`
	IP    float64 `json:"i+"`
	IN    float64 `json:"i-"`
	UXJA  float64 `json:"uxja"`
	UXJB  float64 `json:"uxjb"`
	UXJC  float64 `json:"uxjc"`
	IXJA  float64 `json:"ixja"`
	IXJB  float64 `json:"ixjb"`
	IXJC  float64 `json:"ixjc"`
	UNB   float64 `json:"unb"`
	INB   float64 `json:"inb"`
	PDM   float64 `json:"pdm"`
	QDM   float64 `json:"qdm"`
	SDM   float64 `json:"sdm"`
	IA    float64 `json:"ia"`
	IB    float64 `json:"ib"`
	IC    float64 `json:"ic"`
	UA    float64 `json:"ua"`
	UB    float64 `json:"ub"`
	UC    float64 `json:"uc"`
	PA    float64 `json:"pa"`
	PB    float64 `json:"pb"`
	PC    float64 `json:"pc"`
	QA    float64 `json:"qa"`
	QB    float64 `json:"qb"`
	QC    float64 `json:"qc"`
	SA    float64 `json:"sa"`
	SB    float64 `json:"sb"`
	SC    float64 `json:"sc"`
	PFA   float64 `json:"pfa"`
	PFB   float64 `json:"pfb"`
	PFC   float64 `json:"pfc"`
	UAB   float64 `json:"uab"`
	UBC   float64 `json:"ubc"`
	UCA   float64 `json:"uca"`
	ZYGGL float64 `json:"zyggl"`
	ZWGGL float64 `json:"zwggl"`
	ZSZGL float64 `json:"zszgl"`
	ZGLYS float64 `json:"zglys"`
	F     float64 `json:"f"`
	Time  string  `json:"time"`
	Isend string  `json:"isend"`
}

type MinuteMessage1 struct {
	ID        string  `json:"id"`
	ZYGSZ     float64 `json:"zygsz"`
	FYGSZ     float64 `json:"fygsz"`
	ZWGSZ     float64 `json:"zwgsz"`
	FWGSZ     float64 `json:"fwgsz"`
	ZYJSZ     float64 `json:"zyjsz"`
	FYJSZ     float64 `json:"fyjsz"`
	ZYFSZ     float64 `json:"zyfsz"`
	FYFSZ     float64 `json:"fyfsz"`
	ZYP_SZ    float64 `json:"zypsz"`
	FYP_SZ    float64 `json:"fypsz"`
	ZYVSZ     float64 `json:"zyvsz"`
	FYVSZ     float64 `json:"fyvsz"`
	ZYDVSZ    float64 `json:"zydvsz"`
	FYDVSZ    float64 `json:"fydvsz"`
	ZY6SZ     float64 `json:"zy6sz"`
	FY6SZ     float64 `json:"fy6sz"`
	DMPMAX    float64 `json:"dmpmax"`
	DMPMAXOCT string  `json:"dmpmaxoct"`
	DMSMAX    float64 `json:"dmsmax"`
	DMSMAXOCT string  `json:"dmsmaxoct"`
	Time      string  `json:"time"`
	Isend     string  `json:"isend"`
}

type MinuteMessage2 struct {
	ID     string  `json:"id"`
	UATHD  float64 `json:"uathd"`
	UBTHD  float64 `json:"ubthd"`
	UCTHD  float64 `json:"ucthd"`
	IATHD  float64 `json:"iathd"`
	IBTHD  float64 `json:"ibthd"`
	ICTHD  float64 `json:"icthd"`
	UAXBL3 float64 `json:"uaxbl3"`
	UBXBL3 float64 `json:"ubxbl3"`
	UCXBL3 float64 `json:"ucxbl3"`
	IAXBL3 float64 `json:"iaxbl3"`
	IBXBL3 float64 `json:"ibxbl3"`
	ICXBL3 float64 `json:"icxbl3"`
	UAXBL5 float64 `json:"uaxbl5"`
	UBXBL5 float64 `json:"ubxbl5"`
	UCXBL5 float64 `json:"ucxbl5"`
	IAXBL5 float64 `json:"iaxbl5"`
	IBXBL5 float64 `json:"ibxbl5"`
	ICXBL5 float64 `json:"icxbl5"`
	UAXBL7 float64 `json:"uaxbl7"`
	UBXBL7 float64 `json:"ubxbl7"`
	UCXBL7 float64 `json:"ucxbl7"`
	IAXBL7 float64 `json:"iaxbl7"`
	IBXBL7 float64 `json:"ibxbl7"`
	ICXBL7 float64 `json:"icxbl7"`
	Time   string  `json:"time"`
	Isend  string  `json:"isend"`
}

type MinuteMessage3 struct {
	ID    string  `json:"id"`
	IAXB3 float64 `json:"iaxb3"`
	IBXB3 float64 `json:"ibxb3"`
	ICXB3 float64 `json:"icxb3"`
	IAXB5 float64 `json:"iaxb5"`
	IBXB5 float64 `json:"ibxb5"`
	ICXB5 float64 `json:"icxb5"`
	IAXB7 float64 `json:"iaxb7"`
	IBXB7 float64 `json:"ibxb7"`
	ICXB7 float64 `json:"icxb7"`
	Time  string  `json:"time"`
	Isend string  `json:"isend"`
}

type CombinedMinuteMessage struct {
	ID        string  `json:"id"`
	ZYGSZ     float64 `json:"zygsz"`
	FYGSZ     float64 `json:"fygsz"`
	ZWGSZ     float64 `json:"zwgsz"`
	FWGSZ     float64 `json:"fwgsz"`
	ZYJSZ     float64 `json:"zyjsz"`
	FYJSZ     float64 `json:"fyjsz"`
	ZYFSZ     float64 `json:"zyfsz"`
	FYFSZ     float64 `json:"fyfsz"`
	ZYP_SZ    float64 `json:"zypsz"`
	FYP_SZ    float64 `json:"fypsz"`
	ZYVSZ     float64 `json:"zyvsz"`
	FYVSZ     float64 `json:"fyvsz"`
	ZYDVSZ    float64 `json:"zydvsz"`
	FYDVSZ    float64 `json:"fydvsz"`
	ZY6SZ     float64 `json:"zy6sz"`
	FY6SZ     float64 `json:"fy6sz"`
	DMPMAX    float64 `json:"dmpmax"`
	DMPMAXOCT string  `json:"dmpmaxoct"`
	DMSMAX    float64 `json:"dmsmax"`
	DMSMAXOCT string  `json:"dmsmaxoct"`
	UATHD     float64 `json:"uathd"`
	UBTHD     float64 `json:"ubthd"`
	UCTHD     float64 `json:"ucthd"`
	IATHD     float64 `json:"iathd"`
	IBTHD     float64 `json:"ibthd"`
	ICTHD     float64 `json:"icthd"`
	UAXBL3    float64 `json:"uaxbl3"`
	UBXBL3    float64 `json:"ubxbl3"`
	UCXBL3    float64 `json:"ucxbl3"`
	IAXBL3    float64 `json:"iaxbl3"`
	IBXBL3    float64 `json:"ibxbl3"`
	ICXBL3    float64 `json:"icxbl3"`
	UAXBL5    float64 `json:"uaxbl5"`
	UBXBL5    float64 `json:"ubxbl5"`
	UCXBL5    float64 `json:"ucxbl5"`
	IAXBL5    float64 `json:"iaxbl5"`
	IBXBL5    float64 `json:"ibxbl5"`
	ICXBL5    float64 `json:"icxbl5"`
	UAXBL7    float64 `json:"uaxbl7"`
	UBXBL7    float64 `json:"ubxbl7"`
	UCXBL7    float64 `json:"ucxbl7"`
	IAXBL7    float64 `json:"iaxbl7"`
	IBXBL7    float64 `json:"ibxbl7"`
	ICXBL7    float64 `json:"icxbl7"`
	IAXB3     float64 `json:"iaxb3"`
	IBXB3     float64 `json:"ibxb3"`
	ICXB3     float64 `json:"icxb3"`
	IAXB5     float64 `json:"iaxb5"`
	IBXB5     float64 `json:"ibxb5"`
	ICXB5     float64 `json:"icxb5"`
	IAXB7     float64 `json:"iaxb7"`
	IBXB7     float64 `json:"ibxb7"`
	ICXB7     float64 `json:"icxb7"`
	Time      string  `json:"time"`
	Isend     string  `json:"isend"`
}

var (
	messageBuffer       = make(map[string]*MQTTMessage)
	bufferMutex         = &sync.Mutex{}
	minuteMessageBuffer = make(map[string]*CombinedMinuteMessage)
	minuteBufferMutex   = &sync.Mutex{}
)

func (a *application) mqttMessageHandler(client mqtt.Client, msg mqtt.Message) {
	var sensorValue data.SensorValue
	var mqttMsg MQTTMessage
	// fmt.Println("raw data", string(msg.Payload()))
	// 解析消息
	err := json.Unmarshal(msg.Payload(), &mqttMsg)
	if err != nil {
		a.logger.PrintInfo("Error unmarshaling JSON: %v\n", map[string]string{"error": err.Error()})
		return
	}

	bufferMutex.Lock()
	defer bufferMutex.Unlock()

	// 根据 isend 字段判断是第一个包还是第二个包
	if mqttMsg.Isend == "0" {
		// 第一个包，存储到缓冲区
		messageBuffer[mqttMsg.ID] = &mqttMsg
		return
	} else if mqttMsg.Isend == "1" {
		// 第二个包，检查是否有对应的第一个包
		if existing, exists := messageBuffer[mqttMsg.ID]; exists {
			// 合并两个包的数据
			// 将第二个包的数据复制到第一个包中
			mergeMessages(existing, &mqttMsg)

			// 获取设备ID
			key_id, _ := a.models.SensorDevices.GetSensorIDByDeviceID(existing.ID)
			sensorValue.SensorDeviceID = strconv.Itoa(key_id)

			// 将合并后的消息转换为JSON
			combinedJSON, err := json.Marshal(existing)
			if err != nil {
				a.logger.PrintInfo("Error marshaling combined message: %v\n", map[string]string{"error": err.Error()})
				return
			}

			sensorValue.Value = string(combinedJSON)

			err = a.models.SensorValue.Insert(&sensorValue, key_id)
			if err != nil {
				a.logger.PrintInfo("Error inserting sensor value: %v\n", map[string]string{"error": err.Error()})
			}
			// fmt.Println("data", sensorValue.Value)
			// 删除缓冲区中的数据
			delete(messageBuffer, mqttMsg.ID)
		}
	}

}

func mergeMessages(dst *MQTTMessage, src *MQTTMessage) {
	// 直接复制所有字段，不检查是否为零值
	dst.U0 = src.U0
	dst.UP = src.UP
	dst.UN = src.UN
	dst.I0 = src.I0
	dst.IP = src.IP
	dst.IN = src.IN
	dst.UXJA = src.UXJA
	dst.UXJB = src.UXJB
	dst.UXJC = src.UXJC
	dst.IXJA = src.IXJA
	dst.IXJB = src.IXJB
	dst.IXJC = src.IXJC
	dst.UNB = src.UNB
	dst.INB = src.INB
	dst.PDM = src.PDM
	dst.QDM = src.QDM
	dst.SDM = src.SDM
	// dst.IA = src.IA
	// dst.IB = src.IB
	// dst.IC = src.IC
	// dst.UA = src.UA
	// dst.UB = src.UB
	// dst.UC = src.UC
	// dst.PA = src.PA
	// dst.PB = src.PB
	// dst.PC = src.PC
	// dst.QA = src.QA
	// dst.QB = src.QB
	// dst.QC = src.QC
	// dst.SA = src.SA
	// dst.SB = src.SB
	// dst.SC = src.SC
	// dst.PFA = src.PFA
	// dst.PFB = src.PFB
	// dst.PFC = src.PFC
	// dst.UAB = src.UAB
	// dst.UBC = src.UBC
	// dst.UCA = src.UCA
	// dst.ZYGGL = src.ZYGGL
	// dst.ZWGGL = src.ZWGGL
	// dst.ZSZGL = src.ZSZGL
	// dst.ZGLYS = src.ZGLYS
	// dst.F = src.F
	dst.Isend = "1"
}

func (a *application) mqttDayMessageHandler(client mqtt.Client, msg mqtt.Message) {
	var mqttMsg MQTTMessage
	var sensorValue data.SensorValue

	err := json.Unmarshal(msg.Payload(), &mqttMsg)
	if err != nil {
		a.logger.PrintInfo("Error unmarshaling JSON: %v\n", map[string]string{"error": err.Error()})
		return
	}

	// a.logger.PrintInfo("Received MQTT message: %v\n", map[string]string{"msg": string(msg.Payload())})

	key_id, _ := a.models.SensorDevices.GetSensorIDByDeviceID(mqttMsg.ID)
	sensorValue.SensorDeviceID = strconv.Itoa(key_id)
	sensorValue.Value = string(msg.Payload())

	err = a.models.SensorValue.InsertDay(&sensorValue, key_id)
	if err != nil {
		a.logger.PrintInfo("Error inserting sensor value: %v\n", map[string]string{"error": err.Error()})
	}
}

func (a *application) mqttMinuteMessageHandler(client mqtt.Client, msg mqtt.Message) {
	var sensorValue data.SensorValue
	var combinedMsg CombinedMinuteMessage

	// fmt.Println("raw data", string(msg.Payload()))

	// 解析消息
	var msg1 MinuteMessage1
	var msg2 MinuteMessage2
	var msg3 MinuteMessage3

	// 先尝试解析为第三个包，因为它的字段最少且最独特
	err := json.Unmarshal(msg.Payload(), &msg3)
	if err == nil {
		// 检查是否包含第三个包特有的字段
		if msg3.Isend == "1" {
			// 这是第三个包
			minuteBufferMutex.Lock()
			defer minuteBufferMutex.Unlock()
			if existing, exists := minuteMessageBuffer[msg3.ID]; exists {
				// 合并第三个包的数据
				existing.IAXB3 = msg3.IAXB3
				existing.IBXB3 = msg3.IBXB3
				existing.ICXB3 = msg3.ICXB3
				existing.IAXB5 = msg3.IAXB5
				existing.IBXB5 = msg3.IBXB5
				existing.ICXB5 = msg3.ICXB5
				existing.IAXB7 = msg3.IAXB7
				existing.IBXB7 = msg3.IBXB7
				existing.ICXB7 = msg3.ICXB7
				existing.Isend = "1"

				// 检查是否所有包都已收到
				if existing.DMPMAX != 0 && existing.UATHD != 0 {
					// 获取设备ID
					key_id, _ := a.models.SensorDevices.GetSensorIDByDeviceID(existing.ID)
					sensorValue.SensorDeviceID = strconv.Itoa(key_id)

					// 将合并后的消息转换为JSON
					combinedJSON, err := json.Marshal(existing)
					if err != nil {
						a.logger.PrintInfo("Error marshaling combined message: %v\n", map[string]string{"error": err.Error()})
						return
					}

					sensorValue.Value = string(combinedJSON)
					err = a.models.SensorValue.InsertMinute(&sensorValue, key_id)
					if err != nil {
						a.logger.PrintInfo("Error inserting sensor value: %v\n", map[string]string{"error": err.Error()})
					}

					// 删除缓冲区中的数据
					delete(minuteMessageBuffer, msg3.ID)
				}
			}
			return
		}
	}

	// 尝试解析为第一个包，通过检查其特有的字段
	err = json.Unmarshal(msg.Payload(), &msg1)
	if err == nil {
		// 检查是否包含第一个包特有的字段
		if msg1.DMPMAX != 0 || msg1.DMSMAX != 0 {
			// 这是第一个包
			minuteBufferMutex.Lock()
			defer minuteBufferMutex.Unlock()

			if existing, exists := minuteMessageBuffer[msg1.ID]; exists {
				// 合并第一个包的数据
				existing.ZYGSZ = msg1.ZYGSZ
				existing.FYGSZ = msg1.FYGSZ
				existing.ZWGSZ = msg1.ZWGSZ
				existing.FWGSZ = msg1.FWGSZ
				existing.ZYJSZ = msg1.ZYJSZ
				existing.FYJSZ = msg1.FYJSZ
				existing.ZYFSZ = msg1.ZYFSZ
				existing.FYFSZ = msg1.FYFSZ
				existing.ZYP_SZ = msg1.ZYP_SZ
				existing.FYP_SZ = msg1.FYP_SZ
				existing.ZYVSZ = msg1.ZYVSZ
				existing.FYVSZ = msg1.FYVSZ
				existing.ZYDVSZ = msg1.ZYDVSZ
				existing.FYDVSZ = msg1.FYDVSZ
				existing.ZY6SZ = msg1.ZY6SZ
				existing.FY6SZ = msg1.FY6SZ
				existing.DMPMAX = msg1.DMPMAX
				existing.DMPMAXOCT = msg1.DMPMAXOCT
				existing.DMSMAX = msg1.DMSMAX
				existing.DMSMAXOCT = msg1.DMSMAXOCT
				existing.Time = msg1.Time
			} else {
				combinedMsg.ID = msg1.ID
				combinedMsg.Time = msg1.Time
				combinedMsg.ZYGSZ = msg1.ZYGSZ
				combinedMsg.FYGSZ = msg1.FYGSZ
				combinedMsg.ZWGSZ = msg1.ZWGSZ
				combinedMsg.FWGSZ = msg1.FWGSZ
				combinedMsg.ZYJSZ = msg1.ZYJSZ
				combinedMsg.FYJSZ = msg1.FYJSZ
				combinedMsg.ZYFSZ = msg1.ZYFSZ
				combinedMsg.FYFSZ = msg1.FYFSZ
				combinedMsg.ZYP_SZ = msg1.ZYP_SZ
				combinedMsg.FYP_SZ = msg1.FYP_SZ
				combinedMsg.ZYVSZ = msg1.ZYVSZ
				combinedMsg.FYVSZ = msg1.FYVSZ
				combinedMsg.ZYDVSZ = msg1.ZYDVSZ
				combinedMsg.FYDVSZ = msg1.FYDVSZ
				combinedMsg.ZY6SZ = msg1.ZY6SZ
				combinedMsg.FY6SZ = msg1.FY6SZ
				combinedMsg.DMPMAX = msg1.DMPMAX
				combinedMsg.DMPMAXOCT = msg1.DMPMAXOCT
				combinedMsg.DMSMAX = msg1.DMSMAX
				combinedMsg.DMSMAXOCT = msg1.DMSMAXOCT
				minuteMessageBuffer[msg1.ID] = &combinedMsg
			}
			return
		}
	}

	// 尝试解析为第二个包，通过检查其特有的字段
	err = json.Unmarshal(msg.Payload(), &msg2)
	if err == nil {
		// 检查是否包含第二个包特有的字段
		if msg2.UATHD != 0 || msg2.UBTHD != 0 || msg2.UCTHD != 0 {
			// 这是第二个包
			minuteBufferMutex.Lock()
			defer minuteBufferMutex.Unlock()

			if existing, exists := minuteMessageBuffer[msg2.ID]; exists {
				// 合并第二个包的数据
				existing.UATHD = msg2.UATHD
				existing.UBTHD = msg2.UBTHD
				existing.UCTHD = msg2.UCTHD
				existing.IATHD = msg2.IATHD
				existing.IBTHD = msg2.IBTHD
				existing.ICTHD = msg2.ICTHD
				existing.UAXBL3 = msg2.UAXBL3
				existing.UBXBL3 = msg2.UBXBL3
				existing.UCXBL3 = msg2.UCXBL3
				existing.IAXBL3 = msg2.IAXBL3
				existing.IBXBL3 = msg2.IBXBL3
				existing.ICXBL3 = msg2.ICXBL3
				existing.UAXBL5 = msg2.UAXBL5
				existing.UBXBL5 = msg2.UBXBL5
				existing.UCXBL5 = msg2.UCXBL5
				existing.IAXBL5 = msg2.IAXBL5
				existing.IBXBL5 = msg2.IBXBL5
				existing.ICXBL5 = msg2.ICXBL5
				existing.UAXBL7 = msg2.UAXBL7
				existing.UBXBL7 = msg2.UBXBL7
				existing.UCXBL7 = msg2.UCXBL7
				existing.IAXBL7 = msg2.IAXBL7
				existing.IBXBL7 = msg2.IBXBL7
				existing.ICXBL7 = msg2.ICXBL7
			} else {
				combinedMsg.ID = msg2.ID
				combinedMsg.Time = msg2.Time
				combinedMsg.UATHD = msg2.UATHD
				combinedMsg.UBTHD = msg2.UBTHD
				combinedMsg.UCTHD = msg2.UCTHD
				combinedMsg.IATHD = msg2.IATHD
				combinedMsg.IBTHD = msg2.IBTHD
				combinedMsg.ICTHD = msg2.ICTHD
				combinedMsg.UAXBL3 = msg2.UAXBL3
				combinedMsg.UBXBL3 = msg2.UBXBL3
				combinedMsg.UCXBL3 = msg2.UCXBL3
				combinedMsg.IAXBL3 = msg2.IAXBL3
				combinedMsg.IBXBL3 = msg2.IBXBL3
				combinedMsg.ICXBL3 = msg2.ICXBL3
				combinedMsg.UAXBL5 = msg2.UAXBL5
				combinedMsg.UBXBL5 = msg2.UBXBL5
				combinedMsg.UCXBL5 = msg2.UCXBL5
				combinedMsg.IAXBL5 = msg2.IAXBL5
				combinedMsg.IBXBL5 = msg2.IBXBL5
				combinedMsg.ICXBL5 = msg2.ICXBL5
				combinedMsg.UAXBL7 = msg2.UAXBL7
				combinedMsg.UBXBL7 = msg2.UBXBL7
				combinedMsg.UCXBL7 = msg2.UCXBL7
				combinedMsg.IAXBL7 = msg2.IAXBL7
				combinedMsg.IBXBL7 = msg2.IBXBL7
				combinedMsg.ICXBL7 = msg2.ICXBL7
				minuteMessageBuffer[msg2.ID] = &combinedMsg
			}
			return
		}
	}

	// 如果所有解析都失败
	a.logger.PrintInfo("Error unmarshaling JSON: %v\n", map[string]string{"error": "Failed to parse message"})
}
