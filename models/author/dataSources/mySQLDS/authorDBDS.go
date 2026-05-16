package mySQLDS

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/models/author/dataModel"
	"MyProject/models/author/dataSources"
	"MyProject/pkg/pagination"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type AuthorDBDS struct {
	tableName string
	db        *sql.DB
}

func init() {
	validate = validator.New()
}

func NewAuthorDBDS(tableName string, db *sql.DB) (dataSources.AuthorDS, error) {
	ff := &AuthorDBDS{
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

func (ds *AuthorDBDS) CreateAuthor(ctx context.Context, req authorSchema.CreateAuthor) (res dataModel.Author, err error) {
	err = validate.Struct(req)
	if err != nil {
		var validation []string
		if castValidate, ok := err.(validator.ValidationErrors); ok {
			for _, v := range castValidate {
				validation = append(validation, switchValidateErr(v))
			}
			return dataModel.Author{}, fmt.Errorf("%v", validation)

		} else {
			validation = append(validation, err.Error())
		}
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

func switchValidateErr(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return err.Field() + "this field is required"
	case "max":
		return err.Field() + "this field cannot be more than " + err.Param() + " allowed"
	case "min":
		return err.Field() + "this field cannot be lower than " + err.Param() + "allowed"
	case "email":
		return err.Field() + "this field must be a valid email address"
	case "len":
		return err.Field() + "this field must be " + err.Param() + "."
	default:
		return err.Tag() + err.Field()

	}
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
