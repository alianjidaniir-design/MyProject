package mySQLDS

import (
	"MyProject/apiSchema/bookSchema"
	"MyProject/models/book/dataModel"
	"MyProject/models/book/dataSources"
	"MyProject/models/payment/dataModels"
	"MyProject/pkg/pagination"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type BookDBDS struct {
	tableName string
	db        *sql.DB
}

func NewBookDBDS(tableName string, db *sql.DB) (dataSources.BookDS, error) {
	ff := &BookDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *BookDBDS) RegisterBook(ctx context.Context, req bookSchema.RegistrationBook) (res dataModel.Book, err error) {
	if req.Name == "" {
		return dataModel.Book{}, errors.New("the name is empty")
	}
	if req.Writer == "" {
		return dataModel.Book{}, errors.New("the writer is empty")
	}
	if req.Publisher == "" {
		return dataModel.Book{}, errors.New("the publisher is empty")
	}
	translator := req.Translator
	insertQuery := fmt.Sprintf("INSERT INTO %s (code ,name , writer , translator , publisher ) VALUES (?,? , ? ,? , ?)", ds.tableName)
	_, err = ds.db.ExecContext(ctx, insertQuery, req.Code, req.Name, req.Writer, translator, req.Publisher)
	if err != nil {
		return dataModel.Book{}, err
	}
	return ds.selectBook(ctx, req.Code)

}

func (ds *BookDBDS) DeleteBook(ctx context.Context, req bookSchema.GetCodeBook) (res dataModel.Book, err error) {
	err = ds.checkID(ctx, req.Code)
	if err != nil {
		return dataModel.Book{}, err
	}
	deleted := fmt.Sprintf("DELETE FROM %s WHERE code = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleted, req.Code)
	if err != nil {
		return dataModel.Book{}, err
	}
	return dataModel.Book{}, nil
}

func (ds *BookDBDS) DetailBook(ctx context.Context, req bookSchema.GetCodeBook) (res dataModel.Book, err error) {
	err = ds.checkID(ctx, req.Code)
	if err != nil {
		return dataModel.Book{}, err
	}
	return ds.selectBook(ctx, req.Code)
}
func (ds *BookDBDS) ListBooks(ctx context.Context, req bookSchema.PaginationBook) (res []dataModel.Book, total int, err error) {
	var books []dataModel.Book
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	limit := pageSize
	offset := (page - 1) * limit
	var tot int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&tot)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ?", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var book dataModel.Book
		var translator dataModels.NullString
		err = rows.Scan(&book.Code, &book.Name, &book.Writer, &translator.Val, &book.Publisher)
		if err != nil {
			return nil, 0, err
		}
		book.Translator = translator
		books = append(books, book)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return books, tot, nil
}

func (ds *BookDBDS) selectBook(ctx context.Context, Code int64) (book dataModel.Book, err error) {
	var myBook dataModel.Book
	var translator dataModels.NullString
	selectQuery := fmt.Sprintf("SELECT code , name , writer , translator , publisher FROM %s WHERE code=?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, selectQuery, Code).Scan(&myBook.Code, &myBook.Name, &myBook.Writer, &translator.Val, &myBook.Publisher)
	if err != nil {
		return dataModel.Book{}, err
	}
	myBook.Translator = translator
	return myBook, nil
}

func (ds *BookDBDS) checkID(ctx context.Context, Code int64) error {
	var check bool
	checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM books WHERE code = ?)THEN 1 ELSE 0 END`
	err := ds.db.QueryRowContext(ctx, checkQuery, Code).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("the code does not exist")
	}
	return nil
}
