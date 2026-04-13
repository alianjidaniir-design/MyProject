package teachers

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/teacherSchema"
	"MyProject/models/teachers/dataSources"
	mysqlDataSource "MyProject/models/teachers/dataSources/mysqlDS"
	"MyProject/statics/constants/status"
	"context"
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

func (repo *Repository) db() dataSources.TeacherDS {
	return repo.DBDS
}
