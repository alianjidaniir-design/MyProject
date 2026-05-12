package dataSources

import (
	"MyProject/apiSchema/roleSchema"
	"MyProject/models/role/dataModel"
	"context"
)

type RoleDS interface {
	CreateRole(ctx context.Context, req roleSchema.CreateRoleRequest) (res dataModel.Role, err error)
}
