package course

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/courseSchema"
	courseDataSources "MyProject/models/course/dataSources"
	"MyProject/models/course/dataSources/mySqlDS"
	"MyProject/statics/constants/status"
	"context"
	"fmt"
	"log"
	"sync"
)

type Repository struct {
	DBDS    courseDataSources.CourseDB
	initErr error
}

var (
	once sync.Once
	repo *Repository
)

func instance() {
	load, err := mySqlDS.LoadConfig()
	if err != nil {
		repo = &Repository{initErr: fmt.Errorf("LoadConfig() error", err)}
		return
	}
	dbconn, err := mySqlDS.Open(load)
	if err != nil {
		repo = &Repository{initErr: fmt.Errorf("Open() error", err)}
		return
	}

	instance, err := mySqlDS.NewCourseDBDS(load.CourseTableName, dbconn)
	if err != nil {
		repo = &Repository{initErr: fmt.Errorf("NewCourseDBDS() error", err)}
	}
	repo = &Repository{DBDS: instance}
	log.Printf("Successfully instantiated CourseRepo")

}

func GetRepoIns() *Repository {
	once.Do(instance)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[courseSchema.RequestCourse]) (res courseSchema.ResponseCourse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return courseSchema.ResponseCourse{}, "10", status.UnAvailableServiceError, fmt.Errorf("initErr", repo.initErr)
	}
	if repo.DBDS == nil {
		return courseSchema.ResponseCourse{}, "11", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	createsd, err := repo.db().CreateCourse(ctx, req.Body)
	if err != nil {
		return courseSchema.ResponseCourse{}, "12", status.StatusInternalServerError, err
	}
	return courseSchema.ResponseCourse{Course: createsd}, "0", status.StatusOK, nil
}

func (repo *Repository) db() courseDataSources.CourseDB {
	return repo.DBDS
}
