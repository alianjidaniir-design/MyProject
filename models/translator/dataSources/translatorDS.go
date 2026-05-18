package dataSources

import (
	"MyProject/apiSchema/translatorSchema"
	"MyProject/models/translator/dataModel"
	"context"
)

type TranslatorDS interface {
	CreateTranslator(ctx context.Context, req translateSchema.CreateTranslator) (res dataModel.Translator, err error)
	GetTranslatorAuthor(ctx context.Context, req translateSchema.GetTranslator) (dataModel.Translator, error)
	DeleteTranslator(ctx context.Context, req translateSchema.GetTranslator) (res dataModel.Translator, err error)
	ListTranslator(ctx context.Context, req translateSchema.Pagination) (res []dataModel.Translator, total int, err error)
}
