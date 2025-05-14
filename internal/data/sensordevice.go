package data

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type SensorDevice struct {
	ID           int64  `json:"id"`
	BuildingID   int64  `json:"building_id"`
	DeviceID     string `json:"device_id"`
	SensorTypeID int64  `json:"sensor_type"`
	Name         string `json:"name"`
	Location     string `json:"location"`
}

type SensorDeviceModel struct {
	DB *sql.DB
}

type SensorDeviceWithTypeName struct {
	ID             int64  `json:"id"`
	BuildingID     int64  `json:"building_id"`
	SensorTypeID   int64  `json:"sensor_type"`
	SensorTypeName string `json:"sensor_type_name"`
	Name           string `json:"name"`
	Location       string `json:"location"`
}

func (m SensorDeviceModel) Insert(sensor *SensorDevice) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO sensor_devices (building_id, device_id,sensor_type_id, name, location)
	VALUES (?,?,?,?,?);`
	args := []interface{}{sensor.BuildingID, sensor.DeviceID, sensor.SensorTypeID, sensor.Name, sensor.Location}
	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	sensor.ID = id
	return id, nil
}

func (m SensorDeviceModel) GetAllForBuilding(buildingID int64) ([]SensorDeviceWithTypeName, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	SELECT sd.id, sd.building_id, sd.sensor_type_id, st.name AS sensor_type_name,
	sd.name AS sensor_name, sd.location 
	FROM sensor_devices sd
	JOIN sensor_types st ON sd.sensor_type_id = st.id
	WHERE sd.building_id = ?;`
	rows, err := m.DB.QueryContext(ctx, query, buildingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []SensorDeviceWithTypeName

	for rows.Next() {
		var s SensorDeviceWithTypeName
		err := rows.Scan(
			&s.ID,
			&s.BuildingID,
			&s.SensorTypeID,
			&s.SensorTypeName,
			&s.Name,
			&s.Location,
		)
		if err != nil {
			return nil, err
		}
		sensors = append(sensors, s)
	}
	return sensors, nil
}

// Get the primary key ID by the actual sensor ID
func (m SensorDeviceModel) GetSensorIDByDeviceID(deviceID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	SELECT id FROM sensor_devices WHERE device_id = ?;`
	var id int
	err := m.DB.QueryRowContext(ctx, query, deviceID).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

// Add one sensor device to [sensor_devices] table and create one table for this sensor
func (m SensorDeviceModel) InsertWithTable(sensor *SensorDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Begin transaction
	tx, err := m.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert the sensor device
	query := `
	INSERT INTO sensor_devices (building_id, device_id, sensor_type_id, name, location)
	VALUES (?,?,?,?,?);`
	args := []interface{}{sensor.BuildingID, sensor.DeviceID, sensor.SensorTypeID, sensor.Name, sensor.Location}
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sensor.ID = id

	createTableQuery := fmt.Sprintf(`
	CREATE TABLE IF NOT EXISTS sensor_value_%d (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		sensor_device_id VARCHAR(255) NOT NULL,
		value text NOT NULL COMMENT 'Sensor reading',
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS sensor_value_%d_day (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		sensor_device_id VARCHAR(255) NOT NULL,
		value text NOT NULL COMMENT 'Sensor reading',
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS sensor_value_%d_minute (
		id BIGINT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
		sensor_device_id VARCHAR(255) NOT NULL,
		value text NOT NULL COMMENT 'Sensor reading',
		timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`, sensor.ID, sensor.ID, sensor.ID)

	_, err = tx.ExecContext(ctx, createTableQuery)
	if err != nil {
		return err
	}

	return tx.Commit()
}
