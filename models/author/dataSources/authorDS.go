package dataSources

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/models/author/dataModel"
	"context"
)

type AuthorDS interface {
	CreateAuthor(ctx context.Context, req authorSchema.CreateAuthor) (res dataModel.Author, err error)
}
