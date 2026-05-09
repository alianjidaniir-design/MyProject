package mySQLDS

import (
	"MyProject/apiSchema/categorySchema"
	"MyProject/models/category/dataModel"
	"MyProject/models/category/dataSources"
	"MyProject/pkg/pagination"
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

func (ff *CategoryDBDS) GetDetailCategory(ctx context.Context, req categorySchema.GetRowCategoryRequest) (res dataModel.Category, err error) {
	err = ff.checkCategory(ctx, req.Row)
	if err != nil {
		return dataModel.Category{}, err
	}
	return ff.selectCategory(ctx, req.Row)
}

func (ff *CategoryDBDS) ListCategory(ctx context.Context, req categorySchema.PaginationList) (res []dataModel.Category, total int, err error) {
	var categories []dataModel.Category
	page, pageSize, err := pagination.CheckPage(req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	limit := pageSize
	offset := (page - 1) * limit
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ff.tableName)
	err = ff.db.QueryRowContext(ctx, countQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ", ff.tableName)
	rows, err := ff.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var category dataModel.Category
		err = rows.Scan(&category.Row, &category.Name)
		if err != nil {
			return nil, 0, err
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return categories, totalRows, nil

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
