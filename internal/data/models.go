package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Users            UserModel
	Tokens           TokenModel
	UserGroups       UserGroupModel
	UserGroupMembers UserGroupMembersModel
	Buildings        BuildingModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:            UserModel{DB: db},
		Tokens:           TokenModel{DB: db},
		UserGroups:       UserGroupModel{DB: db},
		UserGroupMembers: UserGroupMembersModel{DB: db},
		Buildings:        BuildingModel{DB: db},
	}
}
