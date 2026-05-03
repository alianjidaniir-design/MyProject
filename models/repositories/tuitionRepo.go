package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition"
	"context"
)

type TuitionRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.CreateTuition]) (res tuitionSchema.InformationTuitionSchema, errStr string, code int, err error)
	Update(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.UpdateTuition]) (res tuitionSchema.MassageTuition, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.DeleteTuition]) (res tuitionSchema.MassageTuition, errStr string, code int, err error)
	ListFixTuition(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.ListFixedTuition]) (res tuitionSchema.ListTuitionSchema, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.ListFixedTuition]) (res tuitionSchema.ListAllTuitionSchema, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.GetTuition]) (res tuitionSchema.TuitionStudentSchema, errStr string, code int, err error)
}

var TuitionRepo TuitionRepository = tuition.GetRepo()
