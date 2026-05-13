package authz

import (
	"MyProject/models/role/dataModel"
	"database/sql"
)

func ParseRole(s string) (*dataModel.Role, error) {
	var db *sql.DB
	role, err := GetRoleByID(db, s)
	if err != nil {
		return nil, err
	}
	return role, nil

}
