package dataSources

import (
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition/dataModels"
	"context"
)

type TuitionDS interface {
	CreateTuition(ctx context.Context, req tuitionSchema.CreateTuition) (res dataModels.Tuition, err error)
	UpdateTuition(ctx context.Context, req tuitionSchema.UpdateTuition) (res dataModels.Tuition, err error)
	DeleteTuition(ctx context.Context, req tuitionSchema.DeleteTuition) (res dataModels.Tuition, err error)
	ListFixedTuition(ctx context.Context, req tuitionSchema.ListFixedTuition) (res []dataModels.StudentsDebit, err error, total int)
	ListAllTuitionStudents(ctx context.Context, req tuitionSchema.ListFixedTuition) (res []dataModels.Tuition, err error, total int)
	GetTuitionStudent(ctx context.Context, req tuitionSchema.GetTuition) (res []dataModels.TuitionStudent, Units int, debits int, credits int, reminder int, err error)
}
