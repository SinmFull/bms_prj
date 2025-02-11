package data

import (
	"context"
	"database/sql"
	"time"
)

type UserGroup struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type UserGroupModel struct {
	DB *sql.DB
}

func (m UserGroupModel) Create(ug *UserGroup) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	INSERT INTO user_groups (name)
	VALUES (?);`

	result, err := m.DB.ExecContext(ctx, query, ug.Name)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ug.ID = id
	return nil
}

func (m UserGroupModel) Get(name string) (*UserGroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	SELECT id, name
	FROM user_groups
	WHERE name = ?;`

	var ug UserGroup
	err := m.DB.QueryRowContext(ctx, query, name).Scan(&ug.ID, &ug.Name)
	if err != nil {
		return nil, err
	}
	return &ug, nil
}
