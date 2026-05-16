package mySQLDS

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/models/author/dataModel"
	"MyProject/models/author/dataSources"
	"context"
	"database/sql"
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
	now := time.Now().In(myLocation())
	insertQuery := fmt.Sprintf("INSERT INTO %s (first_name , last_name, birth_year , created_at , updated_at) VALUES (?,?,? ,? ,? )", ds.tableName)
	result, err := ds.db.ExecContext(ctx, insertQuery, req.FirstName, req.LastName, req.BirthYear, now, now)
	if err != nil {
		return dataModel.Author{}, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		return dataModel.Author{}, err
	}
	return ds.SelectAuthor(ctx, lastId)

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
	err = ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&author.ID, &author.FirstName, &author.LastName, &author.BirthYear, &author.CreatedAt, &author.UpdatedAt, &author.DeletedAt)
	if err != nil {
		return dataModel.Author{}, err
	}
	return author, nil
}
