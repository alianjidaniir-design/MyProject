package dataSources

import (
	"MyProject/apiSchema/profileSchema"
	"MyProject/models/profile/dataModels"
	"context"
)

type ProfileDS interface {
	CreateScoreStudent(ctx context.Context, req profileSchema.CreateScoresReq) (res dataModels.Profile, err error)
	ListScoresStudents(crx context.Context, req profileSchema.ListAllScoresReq) (res []dataModels.ScoresStudents, total int, err error)
	ListSummeryStudents(ctx context.Context, req profileSchema.ListAllScoresReq) (res []dataModels.StudentsSummary, total int, err error)
	GetStudent(ctx context.Context, req profileSchema.GetScoresReq) (res []dataModels.ScoresAnnouncement, err error)
	DeleteProfile(ctx context.Context, req profileSchema.DeleteScoresReq) (err error)
}
