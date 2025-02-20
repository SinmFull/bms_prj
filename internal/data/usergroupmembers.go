package data

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var ErrInvalidRole = errors.New("invalid role")
var ErrAlreadyExists = errors.New("record already exists")

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

	var exists bool
	checkQuery := `
	SELECT EXISTS(
		SELECT 1 FROM user_group_members WHERE user_id = ? AND group_id = ?
	);`
	err := m.DB.QueryRowContext(ctx, checkQuery, ugm.UserID, ugm.GroupId).Scan(&exists)
	if err != nil {
		return err
	}

	if exists {
		return ErrAlreadyExists
	}

	query := `
	INSERT INTO user_group_members (user_id,group_id,role)
	VALUES (?,?,?);`

	_, err = m.DB.ExecContext(ctx, query, ugm.UserID, ugm.GroupId, ugm.Role)
	if err != nil {
		return err
	}
	return nil
}

func (m UserGroupMembersModel) GetMembers(usg UserGroup) ([]User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	groupID := usg.ID
	query := `
	SELECT user_id
	FROM user_group_members WHERE group_id = ?;`
	rows, err := m.DB.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var usersID []int64
	for rows.Next() {
		var uID int64
		if err := rows.Scan(&uID); err != nil {
			return nil, err
		}
		usersID = append(usersID, uID)
	}
	var users []User
	for _, uID := range usersID {
		var u User
		query := `
		SELECT id, created_at, name, email, role
		FROM users  WHERE id = ?;`
		err := m.DB.QueryRowContext(ctx, query, uID).Scan(
			&u.ID,
			&u.CreatedAt,
			&u.Name,
			&u.Email,
			&u.Role,
		)
		if err != nil {
			return nil, err
		}
		if u.Role != "Admin" {
			users = append(users, u)
		}
	}
	return users, nil
}
