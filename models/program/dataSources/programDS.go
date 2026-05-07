package dataSources

import (
	"MyProject/apiSchema/programSchema"
	"MyProject/models/program/dataModel"
	"context"
)

type ProgramDS interface {
	CreateProgram(ctx context.Context, req programSchema.CreateProgramRequest) (res dataModel.Program, err error)
	GetProgram(ctx context.Context, req programSchema.GetDetailProgramRequest) (res dataModel.Program, err error)
}
