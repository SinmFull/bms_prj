package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SensorValue struct {
	ID             int64     `json:"-"`
	SensorDeviceID string    `json:"sensor_device_id"`
	Value          string    `json:"value"`
	RecordedAt     time.Time `json:"recorded_at"`
}

type SensorVaueModel struct {
	DB *sql.DB
}

func (m SensorVaueModel) Insert(sensorValue *SensorValue, key_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tableName := fmt.Sprintf("sensor_value_%d", key_id)

	query := fmt.Sprintf("INSERT INTO %s (sensor_device_id, value) VALUES (?,?);", tableName)

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

func (m SensorVaueModel) InsertDay(sensorValue *SensorValue, key_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tableName := fmt.Sprintf("sensor_value_%d_day", key_id)

	query := fmt.Sprintf("INSERT INTO %s (sensor_device_id, value) VALUES (?,?);", tableName)
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

func (m SensorVaueModel) InsertMinute(sensorValue *SensorValue, key_id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tableName := fmt.Sprintf("sensor_value_%d_minute", key_id)

	query := fmt.Sprintf("INSERT INTO %s (sensor_device_id, value) VALUES (?,?);", tableName)
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

func (m SensorVaueModel) GetNowForDevice(sensorDeviceID int64, tp string) ([]SensorValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var tb_name string
	if tp == "" {
		tb_name = fmt.Sprintf("sensor_value_%d", sensorDeviceID)
	} else {
		tb_name = fmt.Sprintf("sensor_value_%d_%s", sensorDeviceID, tp)
	}
	var query string
	if tp == "" {
		query = fmt.Sprintf(`
		SELECT * FROM %s ORDER BY timestamp DESC LIMIT 2;`, tb_name)
	} else {
		query = fmt.Sprintf(`
		SELECT * FROM %s ORDER BY timestamp DESC LIMIT 1;`, tb_name)
	}

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return []SensorValue{}, err
	}
	defer rows.Close()

	var sensorValues []SensorValue
	for rows.Next() {
		var sensorValue SensorValue
		err := rows.Scan(&sensorValue.ID, &sensorValue.SensorDeviceID, &sensorValue.Value, &sensorValue.RecordedAt)
		if err != nil {
			return []SensorValue{}, err
		}
		sensorValues = append(sensorValues, sensorValue)
	}
	return sensorValues, nil
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

func (m SensorVaueModel) GetBetweenTime(startTime time.Time, endTime time.Time, key_id int) ([]SensorValue, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf("SELECT * FROM sensor_value_%d WHERE timestamp BETWEEN ? AND ?", key_id)
	rows, err := m.DB.QueryContext(ctx, query, startTime, endTime)
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
