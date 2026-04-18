package dataSources

import (
	"MyProject/apiSchema/departmentSchema"
	"MyProject/models/department/dataModels"
	"context"
)

type DepartmentDB interface {
	CreateDepartment(ctx context.Context, req departmentSchema.CreateDepartmentReq) (dataModels.Department, error)
}
