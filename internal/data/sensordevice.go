package data

import (
	"context"
	"database/sql"
	"time"
)

type SensorDevice struct {
	ID           int64  `json:"id"`
	BuildingID   int64  `json:"building_id"`
	SensorTypeID int64  `json:"sensor_type"`
	Name         string `json:"name"`
	Location     string `json:"location"`
}

type SensorDeviceModel struct {
	DB *sql.DB
}

func (m SensorDeviceModel) Insert(sensor *SensorDevice) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO sensor_devices (building_id, sensor_type_id, name, location)
	VALUES (?,?,?,?);`
	args := []interface{}{sensor.BuildingID, sensor.SensorTypeID, sensor.Name, sensor.Location}
	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sensor.ID = id
	return nil
}

func (m SensorDeviceModel) GetAllForBuilding(buildingID int64) ([]SensorDevice, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	SELECT id, sensor_type_id, name, location
	FROM sensor_devices;`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sensors []SensorDevice
	for rows.Next() {
		var s SensorDevice
		err := rows.Scan(
			&s.ID,
			&s.SensorTypeID,
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
