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
	"strconv"
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
	if req.Name == "" {
		return dataModel.Book{}, errors.New("the name is empty")
	}
	var checkAuthor, checkPublisher, checkSubject bool
	checking := `
SELECT EXISTS (SELECT 1 FROM authors WHERE ID=?)
SELECT EXISTS (SELECT 1 FROM publishers WHERE ID=?)
SELECT EXISTS (SELECT 1 FROM subjects WHERE ID=?)`
	err = ds.db.QueryRowContext(ctx, checking, req.ID).Scan(&checkAuthor, &checkPublisher, &checkSubject)
	if err != nil {
		return dataModel.Book{}, err
	}
	if !checkAuthor {
		return dataModel.Book{}, errors.New("there is not author")
	} else if !checkPublisher {
		return dataModel.Book{}, errors.New("there is not publisher")
	} else if !checkSubject {
		return dataModel.Book{}, errors.New("there is not subject")
	}
	measure := strconv.Itoa(req.PublicationYear)
	if len(measure) != 4 {
		return dataModel.Book{}, errors.New("the measure is invalid")
	}
	translator := req.Translator
	insertQuery := fmt.Sprintf("INSERT INTO %s (ID ,name , author_id , translator , publisher_id , publication_year , pages , edition , subject_id , created_at , updated_at ) VALUES (?,? , ? ,? , ? , ? , ? , ? , ? , ? , ?)", ds.tableName)
	now := time.Now().In(myLocation())
	result, err := ds.db.ExecContext(ctx, insertQuery, req.ID, req.Name, req.AuthorID, translator, req.PublisherID, req.PublicationYear, req.Pages, req.Editions, req.SubjectID, now, now)
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
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Book{}, err
	}
	deleted := fmt.Sprintf("DELETE FROM %s WHERE ID = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleted, req.ID)
	if err != nil {
		return dataModel.Book{}, err
	}
	return dataModel.Book{}, nil
}

func (ds *BookDBDS) DetailBook(ctx context.Context, req bookSchema.GetCodeBook) (res dataModel.Book, err error) {
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Book{}, err
	}
	return ds.selectBook(ctx, req.ID)
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
		err = rows.Scan(&book.ID, &book.Name, &book.AuthorId, &translator.Val, &book.PublisherID, &book.PublicationYear, &book.Pages, &book.Edition, &book.SubjectID, &book.CreatedAt, &book.UpdatedAt, &book.DeletedAt)
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

func (ds *BookDBDS) selectBook(ctx context.Context, ID int64) (book dataModel.Book, err error) {
	var myBook dataModel.Book
	var translator dataModels.NullString
	selectQuery := fmt.Sprintf("SELECT ID , name , author_id , translator , publisher_id , publication_year , pages , edition , subject_id , created_at , updated_at , deleted_at FROM %s WHERE code=?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&myBook.ID, &myBook.Name, &myBook.AuthorId, &translator.Val, &myBook.PublisherID, &myBook.PublicationYear, &myBook.Pages, &myBook.Edition, &myBook.SubjectID, &myBook.CreatedAt, &myBook.UpdatedAt, &myBook.DeletedAt)
	if err != nil {
		return dataModel.Book{}, err
	}
	myBook.Translator = translator
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
		return errors.New("the code does not exist")
	}
	return nil
}
