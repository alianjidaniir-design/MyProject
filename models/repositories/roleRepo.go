package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/roleSchema"
	"MyProject/models/role"
	"context"
)

type RoleRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[roleSchema.CreateRoleRequest]) (res roleSchema.DetailRole, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[roleSchema.GetRoleRequest]) (res roleSchema.DetailRole, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[roleSchema.GetRoleRequest]) (res roleSchema.DetailRole, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[roleSchema.Pagination]) (res roleSchema.ListRole, errStr string, code int, err error)
}

var RoleRepo RoleRepository = role.GetRepo()
