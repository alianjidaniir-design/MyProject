package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/programSchema"
	"MyProject/models/program"
	"context"
)

type ProgramRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[programSchema.CreateProgramRequest]) (res programSchema.DetailProgramResponse, errStr string, code int, err error)
}

var ProgramRepo ProgramRepository = program.GetRepo()
