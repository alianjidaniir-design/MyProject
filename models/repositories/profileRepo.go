package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/profileSchema"
	"MyProject/models/profile"
	"context"
)

type ProfileRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[profileSchema.CreateScoresReq]) (res profileSchema.InformationResponse, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[profileSchema.ListAllScoresReq]) (res profileSchema.ListAllScoresResp, errStr string, code int, err error)
	Summery(ctx context.Context, req commonSchema.BaseRequest[profileSchema.ListAllScoresReq]) (res profileSchema.StudentsSummeryResponse, errStr string, code int, err error)
}

var ProfileRepo ProfileRepository = profile.GetRepo()
