package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/departmentSchema"
	"MyProject/models/department"
	"context"
)

type DepartmentRepositories interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.CreateDepartmentReq]) (res departmentSchema.InformationDepartmentResp, errStr string, code int, err error)
}

var DepartmentRepo DepartmentRepositories = department.GetRepo()
