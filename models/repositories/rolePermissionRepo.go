package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/rolePermissionSchema"
	"MyProject/models/rolePermission"

	"context"
)

type RolePermissionRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[rolePermissionSchema.CreateRolePermissionReq]) (res rolePermissionSchema.DetailRolePermission, errStr string, code int, err error)
}

var RolePermissionRepo RolePermissionRepository = rolePermission.GetRepo()
