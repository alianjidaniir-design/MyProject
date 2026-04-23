package dataSources

import (
	"MyProject/apiSchema/profileSchema"
	"MyProject/models/profile/dataModels"
	"context"
)

type ProfileDS interface {
	CreateScoreStudent(ctx context.Context, req profileSchema.CreateScoresReq) (res dataModels.Profile, err error)
	ListScoresStudents(crx context.Context, req profileSchema.ListAllScoresReq) (res dataModels.Profile, err error)
}
