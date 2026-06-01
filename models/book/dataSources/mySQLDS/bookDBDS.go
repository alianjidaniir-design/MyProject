package mySQLDS

import (
	"MyProject/apiSchema/bookSchema"
	"MyProject/models/book/dataModel"
	"MyProject/models/book/dataSources"
	"MyProject/pkg/filter"
	"MyProject/pkg/pagination"
	Val "MyProject/pkg/val"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
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

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		fmt.Println(err)
	}
	return loc
}

func (ds *BookDBDS) RegisterBook(ctx context.Context, req bookSchema.RegistrationBook) (res dataModel.Book, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Book{}, err
	}
	var checkAuthor, checkPublisher, checkSubject, checkTranslator bool
	checking := `
SELECT
    EXISTS (SELECT 1 FROM authors WHERE ID=?),
    EXISTS (SELECT 1 FROM publishers WHERE ID=?),
    EXISTS (SELECT 1 FROM subjects WHERE ID=?),
    EXISTS(SELECT 1 FROM translators WHERE ID = ? )`

	err = ds.db.QueryRowContext(ctx, checking, req.AuthorID, req.PublisherID, req.SubjectID, req.TranslatorID).Scan(&checkAuthor, &checkPublisher, &checkSubject, &checkTranslator)
	if err != nil {
		return dataModel.Book{}, err
	}
	if !checkAuthor {
		return dataModel.Book{}, errors.New("there is not author")
	} else if !checkPublisher {
		return dataModel.Book{}, errors.New("there is not publisher")
	} else if !checkSubject {
		return dataModel.Book{}, errors.New("there is not subject")
	} else if req.TranslatorID == nil {
	} else if !checkTranslator {
		return dataModel.Book{}, errors.New("there is not translator")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (name , author_id , translator_id , publisher_id , publication_year , pages , edition , subject_id , created_at , updated_at ) VALUES (? , ? ,? , ? , ? , ? , ? , ? , ? , ?)", ds.tableName)
	now := time.Now().In(myLocation())
	result, err := ds.db.ExecContext(ctx, insertQuery, req.Name, req.AuthorID, req.TranslatorID, req.PublisherID, req.PublicationYear, req.Pages, req.Editions, req.SubjectID, now, now)
	if err != nil {
		return dataModel.Book{}, err
	}
	lastID, err := result.LastInsertId()
	if err != nil {
		return dataModel.Book{}, err
	}
	return ds.selectBook(ctx, lastID)

}

func (ds *BookDBDS) DeleteBook(ctx context.Context, req bookSchema.GetCodeBook) (res dataModel.Book, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Book{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Book{}, err
	}
	now := time.Now().In(myLocation())
	deleted := fmt.Sprintf("UPDATE %s SET deleted_at = ? , updated_at = ? WHERE ID = ?", ds.tableName)
	rows, err := ds.db.PrepareContext(ctx, deleted)
	if err != nil {
		return dataModel.Book{}, err
	}
	defer rows.Close()
	_, err = rows.ExecContext(ctx, now, now, req.ID)
	if err != nil {
		return dataModel.Book{}, err
	}
	return ds.selectBook(ctx, req.ID)

	return dataModel.Book{}, nil
}

func (ds *BookDBDS) DetailBook(ctx context.Context, req bookSchema.GetCodeBook) (res dataModel.Book, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Book{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Book{}, err
	}
	return ds.selectBook(ctx, req.ID)
}
func (ds *BookDBDS) ListBooks(ctx context.Context, req bookSchema.PaginationBook) (res []dataModel.Book, total int, err error) {
	var rows *sql.Rows
	err = Val.CheckValidation(req)
	if err != nil {
		return res, total, err
	}
	var books []dataModel.Book
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	limit := pageSize
	offset := (page - 1) * limit
	var tot int
	fil := []filter.Filter{
		{Con: "author_id", Value: req.AuthorID},
		{Con: "translator_id", Value: req.TranslatorID},
		{Con: "publisher_id", Value: req.PublisherID},
		{Con: "subject_id", Value: req.SubjectID},
	}

	cond, args := filter.Filtering(fil...)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", ds.tableName, cond)
	if len(args) > 0 {
		err = ds.db.QueryRowContext(ctx, countQuery, args...).Scan(&tot)
		fmt.Println(ds.db.QueryRowContext(ctx, countQuery, args...).Scan(&tot), tot)
	} else {
		err = ds.db.QueryRowContext(ctx, countQuery).Scan(&tot)
	}
	if err != nil {
		return nil, 0, err
	}

	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE %s LIMIT ? OFFSET ?", ds.tableName, cond)
	queryArgs := append(args, limit, offset)
	rows, err = ds.db.QueryContext(ctx, selectQuery, queryArgs...)
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()
	for rows.Next() {
		var book dataModel.Book
		err = rows.Scan(&book.ID, &book.Name, &book.AuthorId, &book.TranslatorID, &book.PublisherID, &book.PublicationYear, &book.Pages, &book.Edition, &book.SubjectID, &book.CreatedAt, &book.UpdatedAt, &book.DeletedAt)
		if err != nil {
			return nil, 0, err
		}
		books = append(books, book)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return books, tot, nil
}

func (ds *BookDBDS) selectBook(ctx context.Context, ID int64) (book dataModel.Book, err error) {
	var myBook dataModel.Book
	selectQuery := fmt.Sprintf("SELECT ID , name , author_id , translator_id , publisher_id , publication_year , pages , edition , subject_id , created_at , updated_at , deleted_at FROM %s WHERE ID=?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&myBook.ID, &myBook.Name, &myBook.AuthorId, &myBook.TranslatorID, &myBook.PublisherID, &myBook.PublicationYear, &myBook.Pages, &myBook.Edition, &myBook.SubjectID, &myBook.CreatedAt, &myBook.UpdatedAt, &myBook.DeletedAt)
	if err != nil {
		return dataModel.Book{}, err
	}
	return myBook, nil
}

func (ds *BookDBDS) checkID(ctx context.Context, ID int64) error {
	var check bool
	checkQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM books WHERE ID = ?)THEN 1 ELSE 0 END`
	err := ds.db.QueryRowContext(ctx, checkQuery, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("the book does not exist")
	}
	return nil
}
