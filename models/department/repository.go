package department

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/departmentSchema"
	"MyProject/models/department/dataSources"
	mySQLDS "MyProject/models/department/dataSources/mySqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.DepartmentDB
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
		repo = &Repository{initRepo: errors.New("can't open DB")}
		return
	}
	newDepartment, err := mySQLDS.NewDepartmentDBDS(cfg.DepartmentTableName, open)
	if err != nil {
		repo = &Repository{initRepo: errors.New("problem")}
		return
	}
	repo = &Repository{DBDS: newDepartment}
	log.Printf("success create department repository")

}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.CreateDepartmentReq]) (res departmentSchema.InformationDepartmentResp, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return departmentSchema.InformationDepartmentResp{}, "10", status.UnAvailableServiceError, fmt.Errorf("initErr", repo.initRepo)
	}
	if repo.DBDS == nil {
		return departmentSchema.InformationDepartmentResp{}, "11", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	createSd, err := repo.db().CreateDepartment(ctx, req.Body)
	if err != nil {
		return departmentSchema.InformationDepartmentResp{}, "12", status.StatusInternalServerError, err
	}
	return departmentSchema.InformationDepartmentResp{Department: createSd}, "0", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.UpdateDepartmentReq]) (res departmentSchema.InformationDepartmentResp, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return departmentSchema.InformationDepartmentResp{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return departmentSchema.InformationDepartmentResp{}, "02", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	update, err := repo.db().UpdateDepartment(ctx, req.Body)
	if err != nil {
		return departmentSchema.InformationDepartmentResp{}, "03", status.StatusInternalServerError, err
	}
	return departmentSchema.InformationDepartmentResp{Department: update}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.ListReq]) (res departmentSchema.ListDepartmentResp, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return departmentSchema.ListDepartmentResp{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return departmentSchema.ListDepartmentResp{}, "02", status.StatusBadRequest, err
	}
	list, totalAll, err := repo.db().ListDepartment(ctx, req.Body)
	if err != nil {
		return departmentSchema.ListDepartmentResp{}, "03", status.StatusInternalServerError, err
	}
	return departmentSchema.ListDepartmentResp{Department: list, Total: totalAll}, "", 0, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[departmentSchema.DeleteDepartmentReq]) (res departmentSchema.DeleteDepartmentResp, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return departmentSchema.DeleteDepartmentResp{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return departmentSchema.DeleteDepartmentResp{}, "02", status.StatusBadRequest, err
	}
	_, err = repo.db().DeleteDepartment(ctx, req.Body)
	if err != nil {
		return departmentSchema.DeleteDepartmentResp{}, "03", status.StatusInternalServerError, err
	}
	return departmentSchema.DeleteDepartmentResp{Message: "deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.DepartmentDB {
	return repo.DBDS
}
