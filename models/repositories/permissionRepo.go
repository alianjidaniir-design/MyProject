package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/permissionSchema"
	"MyProject/models/permission"
	"context"
)

type PermissionRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.CreatePermissionReq]) (res permissionSchema.DetailPermission, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.GetPermissionReq]) (res permissionSchema.DetailPermission, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.GetPermissionReq]) (res permissionSchema.DetailPermission, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.Pagination]) (res permissionSchema.ListPermission, errStr string, code int, err error)
}

var PermissionRepo PermissionRepository = permission.GetRepo()
