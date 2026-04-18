package dataSources

import (
	"MyProject/apiSchema/termSchema"
	"MyProject/models/term/dataModels"
	"context"
)

type TermDS interface {
	CreateTerm(ctx context.Context, req termSchema.CreateTerm) (res dataModels.Term, err error)
	ListTerms(ctx context.Context, req termSchema.ListTerm) (res []dataModels.Term, total int, err error)
}
