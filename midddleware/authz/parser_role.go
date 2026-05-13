package authz

import (
	"MyProject/models/role/dataModel"
	"strings"
)

func ParseRole(s string) (*dataModel.Role, error) {
	towerName := strings.ToLower(s)
	getRole, err := GetRoleByID(towerName)
	if err != nil {
		return nil, err
	}
	return getRole, nil

}
