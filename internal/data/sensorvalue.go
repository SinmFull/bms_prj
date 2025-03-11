package data

import (
	"context"
	"database/sql"
	"time"
)

type SensorValue struct {
	ID             int64     `json:"-"`
	SensorDeviceID int64     `json:"sensor_device_id"`
	Value          string    `json:"value"`
	RecordedAt     time.Time `json:"recorded_at"`
}

type SensorVaueModel struct {
	DB *sql.DB
}

func (m SensorVaueModel) Insert(sensorValue *SensorValue) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO sensor_value (sensor_device_id, value)
	VALUES (?,?);`
	args := []interface{}{sensorValue.SensorDeviceID, sensorValue.Value}
	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sensorValue.ID = id
	return nil
}

func (m SensorVaueModel) GetNowForDevice(sensorDeviceID int64) (SensorValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	SELECT * FROM sensor_value WHERE sensor_device_id = ? ORDER BY timestamp DESC LIMIT 1;`
	row := m.DB.QueryRowContext(ctx, query, sensorDeviceID)
	var sensorValue SensorValue
	err := row.Scan(&sensorValue.ID, &sensorValue.SensorDeviceID, &sensorValue.Value, &sensorValue.RecordedAt)
	if err != nil {
		return SensorValue{}, err
	}
	return sensorValue, nil
}

func (m SensorVaueModel) GetNowForAllDevices() ([]SensorValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT sv.*
		FROM sensor_value sv
		JOIN (
			SELECT sensor_device_id, MAX(timestamp) AS max_timestamp
			FROM sensor_value
			GROUP BY sensor_device_id
		) latest
		ON sv.sensor_device_id = latest.sensor_device_id AND sv.timestamp = latest.max_timestamp;`

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	res := []SensorValue{}
	for rows.Next() {
		var sensorValue SensorValue
		err := rows.Scan(&sensorValue.ID, &sensorValue.SensorDeviceID, &sensorValue.Value, &sensorValue.RecordedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, sensorValue)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return res, nil
}
