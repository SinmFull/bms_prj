package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrSensorTypeExist = errors.New("sensor type already exists")

type SensorType struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Unit string `json:"unit"`
}

type SensorTypeModel struct {
	DB *sql.DB
}

func (m SensorTypeModel) Insert(sensorType *SensorType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var count int
	m.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM sensor_types WHERE name =?", sensorType.Name).Scan(&count)
	if count > 0 {
		return ErrSensorTypeExist
	}
	query := `
	INSERT INTO sensor_types (name, unit)
	VALUES (?,?);`
	args := []interface{}{sensorType.Name, sensorType.Unit}
	result, err := m.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sensorType.ID = id
	return nil
}

func (m SensorTypeModel) GetAll() ([]*SensorType, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	query := `
	SELECT id, name, unit
	FROM sensor_types;`
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sensorTypes := []*SensorType{}
	for rows.Next() {
		var sensorType SensorType
		rows.Scan(&sensorType.ID, &sensorType.Name, &sensorType.Unit)
		sensorTypes = append(sensorTypes, &sensorType)
	}
	return sensorTypes, nil
}
