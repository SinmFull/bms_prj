package data

import (
	"context"
	"database/sql"
	"time"
)

type Building struct {
	ID       int64  `json:"id"`
	GroupID  int64  `json:"-"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

type BuildingModel struct {
	DB *sql.DB
}

func (b BuildingModel) Insert(building *Building) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO buildings (group_id, name, location)
	VALUES (?, ?, ?);`

	args := []interface{}{building.GroupID, building.Name, building.Location}
	result, _ := b.DB.ExecContext(ctx, query, args...)

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	building.ID = id

	return nil
}

func (b BuildingModel) Get(id int64) ([]Building, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT b.id, b.group_id, b.name, b.location
	FROM buildings b
	JOIN user_groups ug ON b.group_id = ug.id
	JOIN user_group_members ugm ON ug.id = ugm.group_id
	WHERE ugm.user_id = ?;`

	args := []interface{}{id}
	rows, err := b.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var buildings []Building
	for rows.Next() {
		var b Building
		if err := rows.Scan(&b.ID, &b.GroupID, &b.Name, &b.Location); err != nil {
			return nil, err
		}
		buildings = append(buildings, b)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return buildings, nil
}
