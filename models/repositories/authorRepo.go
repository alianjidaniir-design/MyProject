package repositories

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/author"
	"context"
)

type AuthorRepository interface {
	CreateAuthor(ctx context.Context, req commonSchema.BaseRequest[authorSchema.CreateAuthor]) (res authorSchema.DetailAuthor, errStr string, code int, err error)
}

var AuthorRepo AuthorRepository = author.GetRepo()
