package dataSources

import (
	"MyProject/apiSchema/rolePermissionSchema"
	"MyProject/models/rolePermission/dataModel"
	"context"
)

type RolePermissionDS interface {
	CreatePermission(ctx context.Context, req rolePermissionSchema.CreateRolePermissionReq) (res dataModel.RolePermission, err error)
}
