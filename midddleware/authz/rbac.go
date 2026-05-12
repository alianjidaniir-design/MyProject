package authz

import (
	"MyProject/models/role/dataModel"
	"MyProject/statics/constants/permissions"
)

func HasPermissionByTID(role *dataModel.Role, permission permissions.Permissions) bool {
	return false
}
