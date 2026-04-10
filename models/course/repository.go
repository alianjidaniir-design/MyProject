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

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[courseSchema.CoursesListRequest]) (res courseSchema.CourseListResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return courseSchema.CourseListResponse{}, "", 0, err
	}
	if repo.DBDS == nil {
		return courseSchema.CourseListResponse{}, "11", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	List, pagination, err := repo.db().ListCourse(ctx, req.Body)
	if err != nil {
		return courseSchema.CourseListResponse{}, "", status.StatusInternalServerError, err
	}
	return courseSchema.CourseListResponse{Courses: List, Total: pagination}, "0", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[courseSchema.GetCoursesRequest]) (res courseSchema.GetCoursesResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return courseSchema.GetCoursesResponse{}, "01", 0, err
	}
	if repo.DBDS == nil {
		return courseSchema.GetCoursesResponse{}, "02", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	get, err := repo.db().GetCourse(ctx, req.Body)
	if err != nil {
		return courseSchema.GetCoursesResponse{}, "03", status.StatusInternalServerError, err
	}
	return courseSchema.GetCoursesResponse{Courses: get}, "0", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[courseSchema.UpdateCourseRequest]) (res courseSchema.UpdateCourseResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return courseSchema.UpdateCourseResponse{}, "01", 0, err
	}
	if repo.DBDS == nil {
		return courseSchema.UpdateCourseResponse{}, "02", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	updatesd, err := repo.db().UpdateCourse(ctx, req.Body)
	if err != nil {
		return courseSchema.UpdateCourseResponse{}, "03", status.StatusInternalServerError, err
	}
	return courseSchema.UpdateCourseResponse{Course: updatesd}, "0", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[courseSchema.HardDeleteCourseRequest]) (res courseSchema.HardDeleteCourseResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return courseSchema.HardDeleteCourseResponse{}, "01", 0, err
	}
	if repo.DBDS == nil {
		return courseSchema.HardDeleteCourseResponse{}, "02", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	deleted, err := repo.db().DeleteCourse(ctx, req.Body)
	if err != nil {
		return courseSchema.HardDeleteCourseResponse{}, "03", status.StatusInternalServerError, err
	}
	return courseSchema.HardDeleteCourseResponse{Course: deleted}, "0", status.StatusOK, nil
}

func (repo *Repository) SoftDelete(ctx context.Context, req commonSchema.BaseRequest[courseSchema.SoftDeleteCourseRequest]) (res courseSchema.SoftDeleteCourseResponse, errStr string, code int, err error) {

	if repo.initErr != nil {
		return courseSchema.SoftDeleteCourseResponse{}, "01", status.UnAvailableServiceError, repo.initErr
	}
	if repo.DBDS == nil {
		return courseSchema.SoftDeleteCourseResponse{}, "02", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	deleted, err := repo.db().SoftDelete(ctx, req.Body)
	if err != nil {
		return courseSchema.SoftDeleteCourseResponse{}, "03", status.StatusInternalServerError, err
	}
	return courseSchema.SoftDeleteCourseResponse{Course: deleted}, "0", status.StatusOK, nil
}

func (repo *Repository) DeActive(ctx context.Context, req commonSchema.BaseRequest[courseSchema.DeactiveCourseRequest]) (res courseSchema.DeactivateCourseResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return courseSchema.DeactivateCourseResponse{}, "01", status.UnAvailableServiceError, err
	}
	if repo.DBDS == nil {
		return courseSchema.DeactivateCourseResponse{}, "02", status.StatusBadRequest, fmt.Errorf("DBDS is nil")
	}
	deactive, err := repo.db().DeactiveCourse(ctx, req.Body)
	if err != nil {
		return courseSchema.DeactivateCourseResponse{}, "03", status.StatusInternalServerError, err
	}
	return courseSchema.DeactivateCourseResponse{Massage: deactive}, "0", status.StatusOK, nil

}

func (repo *Repository) db() courseDataSources.CourseDB {
	return repo.DBDS
}
