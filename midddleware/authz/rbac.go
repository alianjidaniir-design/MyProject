package authz

import (
	"MyProject/statics/constants/permissions"
	"MyProject/statics/constants/roles"
)

var rolePermissions = map[roles.Role]map[permissions.Permissions]bool{
	roles.RoleAdmin: {
		permissions.UserCreate:     true,
		permissions.UserUpdate:     true,
		permissions.UserDelete:     true,
		permissions.UserGet:        true,
		permissions.UserSoftDelete: true,
		permissions.UserList:       true,
	},
	roles.RoleTeacher: {},
}

func HasPermissionByTID(role roles.Role, permission permissions.Permissions) bool {
	perms, ok := rolePermissions[role]
	if !ok {
		return false
	}
	exists := perms[permission]
	return exists
}
