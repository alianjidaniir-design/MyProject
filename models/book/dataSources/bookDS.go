package dataSources

import (
	"MyProject/apiSchema/bookSchema"
	"MyProject/models/book/dataModel"
	"context"
)

type BookDS interface {
	RegisterBook(ctx context.Context, req bookSchema.RegistrationBook) (res dataModel.Book, err error)
	DeleteBook(ctx context.Context, req bookSchema.GetCodeBook) (res dataModel.Book, err error)
	DetailBook(ctx context.Context, req bookSchema.GetCodeBook) (res dataModel.Book, err error)
	ListBooks(ctx context.Context, req bookSchema.PaginationBook) (res []dataModel.Book, total int, err error)
}
