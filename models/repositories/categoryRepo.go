package repositories

import (
	"MyProject/apiSchema/categorySchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/category"
	"context"
)

type CategoryRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[categorySchema.CreateCategoryRequest]) (res categorySchema.InformationCategoryResponse, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[categorySchema.GetRowCategoryRequest]) (res categorySchema.InformationCategoryResponse, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[categorySchema.GetRowCategoryRequest]) (res categorySchema.InformationCategoryResponse, errStr string, code int, err error)
}

var CategoryRepo CategoryRepository = category.GetRepo()
