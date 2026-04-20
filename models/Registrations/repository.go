package enrollmentsdd

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/enrollmentsdd/dataSources"
	"MyProject/models/enrollmentsdd/dataSources/mysqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"fmt"
	"sync"
)

type Repository struct {
	DBDS     dataSources.EnrollmentDS
	initRepo error
}

var (
	onceEnrollment sync.Once
	repo           *Repository
)

func instance() {
	dsn, err := mysqlDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: errors.New("Problem in config")}
		return
	}
	db, err := mysqlDS.Open(dsn)
	if err != nil {
		repo = &Repository{initRepo: errors.New("Problem in opening database connection")}
		return
	}

	newEnr, err := mysqlDS.NewEnrollmentDBDS(dsn.EnrollmentTableName, db)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}

	repo = &Repository{DBDS: newEnr}

}

func GetRepo() *Repository {
	onceEnrollment.Do(instance)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.EnrollmentRequest]) (res registrationSchema.EnrollmentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.EnrollmentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.EnrollmentResponse{}, "02", status.StatusBadRequest, err
	}
	create, err := repo.db().EnrollStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.EnrollmentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.EnrollmentResponse{Enrollment: create}, "", status.StatusOK, nil
}

func (repo *Repository) Cancel(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.CancelEnrollmentRequest]) (res registrationSchema.DeactivateEnrollmentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.DeactivateEnrollmentResponse{}, "01", status.UnAvailableServiceError, err
	}
	if repo.DBDS == nil {
		return registrationSchema.DeactivateEnrollmentResponse{}, "02", status.StatusBadRequest, err
	}
	cancel, err, massage := repo.db().CancelEnrollment(ctx, req.Body)
	if err != nil {
		return registrationSchema.DeactivateEnrollmentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.DeactivateEnrollmentResponse{Enrollment: cancel, Result: massage}, "", status.StatusOK, nil
}

func (repo *Repository) ListEnrollment(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.ListEnrollmentsRequest]) (res registrationSchema.ListEnrollmentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ListEnrollmentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.ListEnrollmentResponse{}, "02", status.StatusBadRequest, err
	}
	list, total, err := repo.db().ListEnrollment(ctx, req.Body)
	if err != nil {
		return registrationSchema.ListEnrollmentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ListEnrollmentResponse{Enrollments: list, TotalCount: total}, "", status.StatusOK, nil
}

func (repo *Repository) ListStudentCourse(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.ListStudentCoursesRequest]) (res registrationSchema.ListStudentCoursesResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ListStudentCoursesResponse{}, "01", status.StatusInternalServerError, err
	}
	if repo.DBDS == nil {
		return registrationSchema.ListStudentCoursesResponse{}, "02", status.StatusBadRequest, err
	}
	list, err := repo.db().ListStudentCourses(ctx, req.Body)
	fmt.Println(list)
	if err != nil {
		return registrationSchema.ListStudentCoursesResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ListStudentCoursesResponse{Enrollments: list}, "", status.StatusOK, nil
}

func (repo *Repository) ListCourseStudent(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.ListCourseStudentsRequest]) (res registrationSchema.ListCourseStudentsResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ListCourseStudentsResponse{}, "01", status.StatusInternalServerError, err
	}
	if repo.DBDS == nil {
		return registrationSchema.ListCourseStudentsResponse{}, "02", status.StatusBadRequest, err
	}
	list, err := repo.db().ListCourseStudents(ctx, req.Body)
	if err != nil {
		return registrationSchema.ListCourseStudentsResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ListCourseStudentsResponse{Course: list}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.EnrollmentDS {
	return repo.DBDS
}
