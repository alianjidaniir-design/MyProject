package dataSources

import (
	"MyProject/apiSchema/permissionSchema"
	"MyProject/models/permission/dataModel"
	"context"
)

type PermissionDS interface {
	CreatePermission(ctx context.Context, req permissionSchema.CreatePermissionReq) (res dataModel.Permission, err error)
	GetPermission(ctx context.Context, req permissionSchema.GetPermissionReq) (res dataModel.Permission, err error)
	DeletePermission(ctx context.Context, req permissionSchema.GetPermissionReq) (res dataModel.Permission, err error)
	ListPermissions(ctx context.Context, req permissionSchema.Pagination) (res []dataModel.Permission, total int, err error)
	ListPerms(roleID int64) ([]dataModel.Permission, error)
}
