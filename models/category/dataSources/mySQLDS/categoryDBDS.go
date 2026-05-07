package mySQLDS

import (
	"MyProject/apiSchema/categorySchema"
	"MyProject/models/category/dataModel"
	"MyProject/models/category/dataSources"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type CategoryDBDS struct {
	tableName string
	db        *sql.DB
}

func NewCategoryDBDS(tableName string, db *sql.DB) (dataSources.CategoryDS, error) {
	ff := &CategoryDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ff *CategoryDBDS) CreateCategory(ctx context.Context, req categorySchema.CreateCategoryRequest) (res dataModel.Category, err error) {
	var lastInsertId int64
	lastIDQuery := fmt.Sprintf("SELECT COALESCE(MAX(row), 0) FROM %s", ff.tableName)
	err = ff.db.QueryRowContext(ctx, lastIDQuery).Scan(&lastInsertId)
	if err != nil {
		return dataModel.Category{}, err
	}
	newID := lastInsertId + 1
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , name) VALUES (?,?) ", ff.tableName)
	_, err = ff.db.ExecContext(ctx, insertQuery, newID, req.Name)
	if err != nil {
		return dataModel.Category{}, errors.New(err.Error())
	}
	return dataModel.Category{}, nil
}

func (ff *CategoryDBDS) DeleteCategory(ctx context.Context, req categorySchema.GetRowCategoryRequest) (res dataModel.Category, err error) {
	err = ff.checkCategory(ctx, req.Row)
	if err != nil {
		return dataModel.Category{}, err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE row =  ?", ff.tableName)
	_, err = ff.db.ExecContext(ctx, deleteQuery, req.Row)
	if err != nil {
		return dataModel.Category{}, err
	}
	return ff.selectCategory(ctx, req.Row)
}

func (ff *CategoryDBDS) selectCategory(ctx context.Context, ID int64) (res dataModel.Category, err error) {
	var category dataModel.Category
	selectQuery := `SELECT row , name FROM ` + ff.tableName + ` where row = ?`
	err = ff.db.QueryRowContext(ctx, selectQuery, ID).Scan(&category.Row, &category.Name)
	if err != nil {
		return category, errors.New(err.Error())
	}
	return category, nil

}

func (ff *CategoryDBDS) checkCategory(ctx context.Context, ID int64) error {
	var check bool
	selectQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM categories WHERE row = ?) THEN 1 ELSE 0 END`
	err := ff.db.QueryRowContext(ctx, selectQuery, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("category does not exist")
	}
	return nil
}
