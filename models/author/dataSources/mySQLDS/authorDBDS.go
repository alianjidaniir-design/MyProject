package mySQLDS

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/models/author/dataModel"
	"MyProject/models/author/dataSources"
	"MyProject/pkg/pagination"
	Val "MyProject/pkg/val"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type AuthorDBDS struct {
	tableName string
	db        *sql.DB
}

func NewAuthorDBDS(tableName string, db *sql.DB) (dataSources.AuthorDS, error) {
	ff := &AuthorDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *AuthorDBDS) CreateAuthor(ctx context.Context, req authorSchema.CreateAuthor) (res dataModel.Author, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Author{}, err
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (first_name , last_name, birth_year ) VALUES (?,?,?)", ds.tableName)
	result, err := ds.db.ExecContext(ctx, insertQuery, req.FirstName, req.LastName, req.BirthYear)
	if err != nil {
		return dataModel.Author{}, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return dataModel.Author{}, err
	}
	return ds.SelectAuthor(ctx, lastId)

}

func (ds *AuthorDBDS) DeleteAuthor(ctx context.Context, req authorSchema.GetAuthor) (res dataModel.Author, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Author{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Author{}, err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return dataModel.Author{}, err
	}
	return dataModel.Author{}, nil

}

func (ds *AuthorDBDS) ListAuthor(ctx context.Context, req authorSchema.Pagination) (res []dataModel.Author, total int, err error) {
	var authors []dataModel.Author
	err = Val.CheckValidation(req)
	if err != nil {
		return nil, 0, err
	}
	page, size, err := pagination.CheckPage(req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	limit := size
	offset := (page - 1) * limit
	var tot int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableName)
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
		var author dataModel.Author
		err = rows.Scan(&author.ID, &author.FirstName, &author.LastName, &author.BirthYear)
		if err != nil {
			return nil, 0, err
		}
		authors = append(authors, author)

	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return authors, tot, nil
}

func (ds *AuthorDBDS) SelectAuthor(ctx context.Context, ID int64) (res dataModel.Author, err error) {
	var author dataModel.Author
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&author.ID, &author.FirstName, &author.LastName, &author.BirthYear)
	if err != nil {
		return dataModel.Author{}, err
	}
	return author, nil
}

func (ds *AuthorDBDS) GetAuthor(ctx context.Context, req authorSchema.GetAuthor) (res dataModel.Author, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModel.Author{}, err
	}
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return dataModel.Author{}, err
	}
	return ds.SelectAuthor(ctx, req.ID)
}

func (ds *AuthorDBDS) checkID(ctx context.Context, ID int64) error {
	var check bool
	checkQuery := `SELECT EXISTS (SELECT 1 FROM ` + ds.tableName + ` where id=?)`
	err := ds.db.QueryRowContext(ctx, checkQuery, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("ID not exists")
	}
	return nil
}
