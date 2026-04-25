package dataSources

import (
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition/dataModels"
	"context"
)

type TuitionDS interface {
	CreateTuition(ctx context.Context, req tuitionSchema.CreateTuition) (res dataModels.Tuition, err error)
}
