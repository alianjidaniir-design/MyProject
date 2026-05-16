package dataSources

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/models/author/dataModel"
	"context"
)

type AuthorDS interface {
	CreateAuthor(ctx context.Context, req authorSchema.CreateAuthor) (res dataModel.Author, err error)
	GetAuthor(ctx context.Context, req authorSchema.GetAuthor) (dataModel.Author, error)
	DeleteAuthor(ctx context.Context, req authorSchema.GetAuthor) (res dataModel.Author, err error)
	ListAuthor(ctx context.Context, req authorSchema.Pagination) (res []dataModel.Author, total int, err error)
}
