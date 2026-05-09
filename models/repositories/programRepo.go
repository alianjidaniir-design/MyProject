package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/programSchema"
	"MyProject/models/program"
	"context"
)

type ProgramRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[programSchema.CreateProgramRequest]) (res programSchema.DetailProgramResponse, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[programSchema.GetDetailProgramRequest]) (res programSchema.DetailProgramResponse, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[programSchema.DeleteProgramRequest]) (res programSchema.DetailProgramResponse, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[programSchema.PaginationListProgramsRequest]) (res programSchema.ListProgramsResponse, errStr string, code int, err error)
}

var ProgramRepo ProgramRepository = program.GetRepo()
