package dataSources

import (
	"MyProject/apiSchema/categorySchema"
	"MyProject/models/category/dataModel"
	"context"
)

type CategoryDS interface {
	CreateCategory(ctx context.Context, req categorySchema.CreateCategoryRequest) (res dataModel.Category, err error)
	DeleteCategory(ctx context.Context, req categorySchema.GetRowCategoryRequest) (res dataModel.Category, err error)
}
