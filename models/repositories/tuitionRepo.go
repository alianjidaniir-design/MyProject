package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition"
	"context"
)

type TuitionRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.CreateTuition]) (res tuitionSchema.InformationTuitionSchema, errStr string, code int, err error)
}

var TuitionRepo TuitionRepository = tuition.GetRepo()
