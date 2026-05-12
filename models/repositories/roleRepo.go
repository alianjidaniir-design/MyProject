package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/roleSchema"
	"context"
)

type RoleRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[roleSchema.CreateRoleRequest]) (res roleSchema.DetailRole, errStr string, code int, err error)
}

var RoleRepo RoleRepository
