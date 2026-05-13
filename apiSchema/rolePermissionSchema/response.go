package rolePermissionSchema

import "MyProject/models/rolePermission/dataModel"

type DetailRolePermission struct {
	Detail  dataModel.RolePermission `json:"detail"`
	Massage string                   `json:"massage"`
}
