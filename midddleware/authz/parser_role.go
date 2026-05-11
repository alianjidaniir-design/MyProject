package authz

import (
	"MyProject/statics/constants/roles"
	"errors"
)

func ParseRole(s string) (roles.Role, error) {
	switch s {
	case string(roles.RoleTeacher):
		return roles.RoleTeacher, nil

	case string(roles.RoleAdmin):
		return roles.RoleAdmin, nil
	case string(roles.RoleStudent):
		return roles.RoleStudent, nil
	default:
		return "", errors.New("invalid role")

	}
}
