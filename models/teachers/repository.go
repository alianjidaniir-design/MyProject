package teachers

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/teacherSchema"
	"MyProject/models/teachers/dataSources"
	mysqlDataSource "MyProject/models/teachers/dataSources/mysqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.TeacherDS
	initRepo error
}

var (
	once    sync.Once
	repoIns *Repository
)

func initIns() {
	DBConn, err := mysqlDataSource.LoadConfig()
	if err != nil {
		return
	}
	open, err := mysqlDataSource.Open(DBConn)
	if err != nil {
		return
	}
	newTeacher, err := mysqlDataSource.NewTeacherDBDS(DBConn.TeacherTableName, open)
	if err != nil {
		repoIns = &Repository{initRepo: fmt.Errorf("error in newTeacherDBDS")}
		return
	}
	repoIns = &Repository{DBDS: newTeacher}
	log.Println("success in newTeacherDBDS")
}

func GetRepo() *Repository {
	once.Do(initIns)
	return repoIns
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.InformationSchema]) (res teacherSchema.TeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.TeacherSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.TeacherSchema{}, "02", status.StatusUnauthorized, err
	}
	create, err := repo.db().CreateTeacher(ctx, req.Body)

	if err != nil {
		return teacherSchema.TeacherSchema{}, "02", status.UnAvailableServiceError, err
	}
	return teacherSchema.TeacherSchema{Teacher: create}, "04", status.StatusOK, err
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.PaginationSchema]) (res teacherSchema.ListSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.ListSchema{}, "01", 0, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.ListSchema{}, "02", status.StatusUnauthorized, errors.New("wrong db connection")
	}
	list, total, err := repo.db().ListTeachers(ctx, req.Body)
	if err != nil {
		return teacherSchema.ListSchema{}, "03", status.UnAvailableServiceError, err
	}
	return teacherSchema.ListSchema{Teachers: list, Total: total}, "04", 0, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.GetTeacherSchema]) (res teacherSchema.DetailTeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.DetailTeacherSchema{}, "01", 0, repo.initRepo
	}
	if repo.DBDS == nil {
		return teacherSchema.DetailTeacherSchema{}, "02", 0, errors.New("wrong db connection")
	}
	get, err := repo.db().GetTeacherById(ctx, req.Body)
	if err != nil {
		return teacherSchema.DetailTeacherSchema{}, "03", status.StatusUnauthorized, err
	}
	return teacherSchema.DetailTeacherSchema{Teacher: get}, "", status.StatusOK, nil

}

func (repo *Repository) HardDelete(ctx context.Context, req commonSchema.BaseRequest[teacherSchema.SelectTeacherSchema]) (res teacherSchema.HardDeleteTeacherSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return teacherSchema.HardDeleteTeacherSchema{}, "01", status.StatusUnauthorized, errors.New("wrong db connection")
	}
	if repo.DBDS == nil {
		return teacherSchema.HardDeleteTeacherSchema{}, "02", status.StatusBadRequest, errors.New("wrong db connection")
	}
	deleted, err := repo.db().HardDeleteTeachers(ctx, req.Body)
	if err != nil {
		return teacherSchema.HardDeleteTeacherSchema{}, "03", status.StatusUnauthorized, err
	}
	return teacherSchema.HardDeleteTeacherSchema{Massage: deleted}, "", status.StatusOK, nil

}

func (repo *Repository) db() dataSources.TeacherDS {
	return repo.DBDS
}
