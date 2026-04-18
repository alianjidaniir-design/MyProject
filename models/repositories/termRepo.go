package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/termSchema"
	"MyProject/models/term"
	"context"
)

type TermRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[termSchema.CreateTerm]) (res termSchema.InformationTerm, errStr string, code int, err error)
}

var TermRepo TermRepository = term.GetRepo()
