package dataSources

import (
	"MyProject/apiSchema/departmentSchema"
	"MyProject/models/department/dataModels"
	"context"
)

type DepartmentDB interface {
	CreateDepartment(ctx context.Context, req departmentSchema.CreateDepartmentReq) (dataModels.Department, error)
	UpdateDepartment(ctx context.Context, req departmentSchema.UpdateDepartmentReq) (dataModels.Department, error)
	ListDepartment(ctx context.Context, req departmentSchema.ListReq) ([]dataModels.Department, int, error)
	DeleteDepartment(ctx context.Context, req departmentSchema.DeleteDepartmentReq) (dataModels.Department, error)
}
