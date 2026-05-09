package dataSources

import (
	"MyProject/apiSchema/programSchema"
	"MyProject/models/program/dataModel"
	"context"
)

type ProgramDS interface {
	CreateProgram(ctx context.Context, req programSchema.CreateProgramRequest) (res dataModel.Program, err error)
	GetProgram(ctx context.Context, req programSchema.GetDetailProgramRequest) (res dataModel.Program, err error)
	DeleteProgram(ctx context.Context, req programSchema.DeleteProgramRequest) (res dataModel.Program, err error)
	ListProgram(ctx context.Context, req programSchema.PaginationListProgramsRequest) (res []dataModel.Program, total int64, err error)
}
