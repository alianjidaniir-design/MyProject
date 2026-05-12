package permissionSchema

import "MyProject/models/permission/dataModel"

type DetailPermission struct {
	Permission dataModel.Permission `json:"permission"`
	Massage    string               `json:"massage"`
}

type ListPermission struct {
	Massage         string                 `json:"massage"`
	Permission      []dataModel.Permission `json:"permissions"`
	TotalPermission int                    `json:"total_permissions"`
}
