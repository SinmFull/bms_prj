package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrInvalidRole = errors.New("invalid role")

type UserGroupMembers struct {
	ID      int64
	UserID  int64
	GroupId int64
	Role    string
}

type UserGroupMembersModel struct {
	DB *sql.DB
}

func (m UserGroupMembersModel) Insert(ugm *UserGroupMembers) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if ugm.Role != "admin" && ugm.Role != "member" {
		return ErrInvalidRole
	}
	query := `
	INSERT INTO user_group_members (user_id,group_id,role)
	VALUES (?,?,?);`

	_, err := m.DB.ExecContext(ctx, query, ugm.UserID, ugm.GroupId, ugm.Role)
	if err != nil {
		return err
	}
	return nil
}
