package category

import (
	"MyProject/apiSchema/categorySchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/category/dataSources"
	"MyProject/models/category/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.CategoryDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepo() {
	cfg, err := mySQLDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	open, err := mySQLDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	newCategory, err := mySQLDS.NewCategoryDBDS(cfg.CategoryTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newCategory}
	log.Printf("successfully init repository")
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[categorySchema.CreateCategoryRequest]) (res categorySchema.InformationCategoryResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return categorySchema.InformationCategoryResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return categorySchema.InformationCategoryResponse{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateCategory(ctx, req.Body)
	if err != nil {
		return categorySchema.InformationCategoryResponse{}, "03", status.StatusBadRequest, err
	}
	return categorySchema.InformationCategoryResponse{Detail: create, Massage: "created successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[categorySchema.GetRowCategoryRequest]) (res categorySchema.InformationCategoryResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return categorySchema.InformationCategoryResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return categorySchema.InformationCategoryResponse{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeleteCategory(ctx, req.Body)
	if err != nil {
		return categorySchema.InformationCategoryResponse{}, "03", status.StatusBadRequest, err
	}
	return categorySchema.InformationCategoryResponse{Detail: deleted, Massage: "deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[categorySchema.GetRowCategoryRequest]) (res categorySchema.InformationCategoryResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return categorySchema.InformationCategoryResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return categorySchema.InformationCategoryResponse{}, "02", status.StatusInternalServerError, err
	}
	get, err := repo.db().GetDetailCategory(ctx, req.Body)
	if err != nil {
		return categorySchema.InformationCategoryResponse{}, "03", status.StatusBadRequest, err
	}
	return categorySchema.InformationCategoryResponse{Detail: get, Massage: "information category"}, "", 0, nil
}

func (repo *Repository) db() dataSources.CategoryDS {
	return repo.DBDS
}
