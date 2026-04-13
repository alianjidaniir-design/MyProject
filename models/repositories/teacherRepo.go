package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/teacherSchema"
	"MyProject/models/teachers"
	"context"
)

type TeacherRepository interface {
	// Create method
	Create(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.InformationSchema]) (res teacherSchema.TeacherSchema, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.PaginationSchema]) (res teacherSchema.ListSchema, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.GetTeacherSchema]) (res teacherSchema.DetailTeacherSchema, errStr string, code int, err error)
	HardDelete(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.SelectTeacherSchema]) (res teacherSchema.HardDeleteTeacherSchema, errStr string, code int, err error)
}

var TeacherRepo TeacherRepository = teachers.GetRepo()
