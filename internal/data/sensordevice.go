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

type SensorDeviceWithTypeName struct {
	ID             int64  `json:"id"`
	BuildingID     int64  `json:"building_id"`
	SensorTypeID   int64  `json:"sensor_type"`
	SensorTypeName string `json:"sensor_type_name"`
	Name           string `json:"name"`
	Location       string `json:"location"`
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
