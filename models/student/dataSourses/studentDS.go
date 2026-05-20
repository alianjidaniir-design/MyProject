package dataSourses

import (
	"MyProject/apiSchema/studentSchema"
	studentDataModel "MyProject/models/student/dataModel"
	"context"
)

type StudentDB interface {
	CreateStudent(ctx context.Context, req studentSchema.SignUpStudent) (studentDataModel.Student, error)
	ReadStudent(ctx context.Context, req studentSchema.ListRequest) ([]studentDataModel.Student, int64, error)
	GetStudent(ctx context.Context, req studentSchema.GetRequest) (studentDataModel.Student, error)
	UpdateStudent(ctx context.Context, req studentSchema.UpdateUserRequest) (studentDataModel.Student, error)
	DeleteStudent(ctx context.Context, req studentSchema.DeleteRequest) (studentDataModel.Student, error)
	SoftDeleteStudent(ctx context.Context, req studentSchema.SoftDeleteRequest) (studentDataModel.Student, error)
	StudentEntry(ctx context.Context, req studentSchema.LoginStudent) (string, error)
}
