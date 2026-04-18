package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/departmentSchema"
	"MyProject/models/department"
	"context"
)

type DepartmentRepositories interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.CreateDepartmentReq]) (res departmentSchema.InformationDepartmentResp, errStr string, code int, err error)
	Update(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.UpdateDepartmentReq]) (res departmentSchema.InformationDepartmentResp, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.ListReq]) (res departmentSchema.ListDepartmentResp, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.DeleteDepartmentReq]) (res departmentSchema.DeleteDepartmentResp, errStr string, code int, err error)
}

var DepartmentRepo DepartmentRepositories = department.GetRepo()
