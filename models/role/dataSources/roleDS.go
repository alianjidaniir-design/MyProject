package dataSources

import (
	"MyProject/apiSchema/roleSchema"
	"MyProject/models/role/dataModel"
	"context"
)

type RoleDS interface {
	CreateRole(ctx context.Context, req roleSchema.CreateRoleRequest) (res dataModel.Role, err error)
	DeleteRole(ctx context.Context, req roleSchema.GetRoleRequest) (res dataModel.Role, err error)
	GetRole(ctx context.Context, req roleSchema.GetRoleRequest) (res dataModel.Role, err error)
	ListRoles(ctx context.Context, req roleSchema.Pagination) (res []dataModel.Role, total int, err error)
}
