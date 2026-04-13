package dataSources

import (
	"MyProject/apiSchema/teacherSchema"
	"MyProject/models/teachers/dataModels"
	"context"
)

type TeacherDS interface {
	CreateTeacher(ctx context.Context, req teacherSchema.InformationSchema) (res dataModels.Teacher, err error)
	ListTeachers(ctx context.Context, req teacherSchema.PaginationSchema) (res []dataModels.Teacher, total int64, err error)
}
