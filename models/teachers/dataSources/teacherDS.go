package dataSources

import (
	"MyProject/apiSchema/teacherSchema"
	"MyProject/models/teachers/dataModels"
	"context"
)

type TeacherDS interface {
	CreateTeacher(ctx context.Context, req teacherSchema.InformationSchema) (res dataModels.Teacher, err error)
	ListTeachers(ctx context.Context, req teacherSchema.PaginationSchema) (res []dataModels.Teacher, total int64, err error)
	GetTeacherById(ctx context.Context, req teacherSchema.GetTeacherSchema) (res dataModels.Teacher, err error)
	HardDeleteTeachers(ctx context.Context, req teacherSchema.SelectTeacherSchema) (res string, err error)
	SoftDeleteTeachers(ctx context.Context, req teacherSchema.SelectTeacherSchema) (res dataModels.Teacher, err error)
	UpdateTeachers(ctx context.Context, req teacherSchema.SelectTeacherSchema) (res dataModels.Teacher, err error)
	LoginTeachers(ctx context.Context, req teacherSchema.LoginTeacherRequest) (access string, refresh string, massage string, err error)
	Refresh(ctx context.Context, req string) (access string, refresh string, err error)
	Logout(ctx context.Context, req teacherSchema.LogoutSchema, refresh string) error
	MyInfo(ctx context.Context, ID int64) (dataModels.InfoTeacher, error)
}
