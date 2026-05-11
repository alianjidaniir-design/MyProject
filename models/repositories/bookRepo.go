package repositories

import (
	"MyProject/apiSchema/bookSchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/book"
	"context"
)

type BookRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[bookSchema.RegistrationBook]) (res bookSchema.InformationBook, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[bookSchema.GetCodeBook]) (res bookSchema.InformationBook, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[bookSchema.GetCodeBook]) (res bookSchema.InformationBook, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[bookSchema.PaginationBook]) (res bookSchema.ListBooks, errStr string, code int, err error)
}

var BookRepo BookRepository = book.GetRepo()
