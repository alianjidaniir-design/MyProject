package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/termSchema"
	"MyProject/models/term"
	"context"
)

type TermRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[termSchema.CreateTerm]) (res termSchema.InformationTerm, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[termSchema.ListTerm]) (res termSchema.ListTerms, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[termSchema.DeleteTerm]) (res termSchema.DeletedTerm, errStr string, code int, err error)
}

var TermRepo TermRepository = term.GetRepo()
