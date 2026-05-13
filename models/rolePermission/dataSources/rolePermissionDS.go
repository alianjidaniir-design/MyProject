package dataSources

import (
	"MyProject/apiSchema/rolePermissionSchema"
	"context"
)

type RolePermissionDS interface {
	CreatePermission(ctx context.Context, req rolePermissionSchema.CreateRolePermissionReq) (err error)
}
