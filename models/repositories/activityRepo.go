package repositories

import (
	"MyProject/apiSchema/activitySchema"
	"MyProject/apiSchema/commonSchema"
	"context"
)

type ActivityRepository interface {
	CreateActivityByTeacher(ctx context.Context, req commonSchema.BaseRequest[activitySchema.CreateActivity]) (res activitySchema.InformationActivity, errStr string, code int, err error)
}

var ActiveRepo ActivityRepository
