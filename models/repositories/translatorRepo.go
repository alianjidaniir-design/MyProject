package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/translatorSchema"
	"MyProject/models/translator"
	"context"
)

type TranslatorRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[translateSchema.CreateTranslator]) (res translateSchema.DetailTranslator, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[translateSchema.GetTranslator]) (res translateSchema.DetailTranslator, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[translateSchema.GetTranslator]) (res translateSchema.DetailTranslator, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[translateSchema.Pagination]) (res translateSchema.ListTranslator, errStr string, code int, err error)
}

var TranslatorRepo TranslatorRepository = translator.GetRepo()
