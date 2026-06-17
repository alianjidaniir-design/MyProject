package dataSources

import (
	"MyProject/apiSchema/activitySchema"
	"MyProject/models/activity/dataModel"
	"context"
)

type ActivityDS interface {
	CreateActivity(ctx context.Context, req activitySchema.CreateActivity, role string, ID int64) (res dataModel.Activity, err error)
}
