package enrollment

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/enrollmentSchema"
	"MyProject/models/enrollment/dataSources"
	"MyProject/models/enrollment/dataSources/mysqlDS"
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
	dsn, err := mysqlDS.LoadConfiger()
	if err != nil {
		repo = &Repository{initRepo: errors.New("Problem in config")}
		return
	}
	db, err := mysqlDS.Open(dsn)
	if err != nil {
		repo = &Repository{initRepo: errors.New("Problem in opening database connection")}
		return
	}

	insta, err := mysqlDS.NewEnrollmentDBDS(dsn.EnrollmentTableName, db)
	if err != nil {
		repo = &Repository{initRepo: err}
	}
	repo = &Repository{DBDS: insta}

}

func GetRepo() *Repository {
	onceEnrollment.Do(instance)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.EnrollmentRequest]) (res enrollmentSchema.EnrollmentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return enrollmentSchema.EnrollmentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return enrollmentSchema.EnrollmentResponse{}, "02", status.StatusBadRequest, err
	}
	create, err := repo.db().EnrollStudent(ctx, req.Body)
	if err != nil {
		return enrollmentSchema.EnrollmentResponse{}, "03", status.StatusInternalServerError, err
	}
	return enrollmentSchema.EnrollmentResponse{Enrollment: create}, "", status.StatusOK, nil
}

func (repo *Repository) Cancel(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.CancelEnrollmentRequest]) (res enrollmentSchema.DeactivateEnrollmentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return enrollmentSchema.DeactivateEnrollmentResponse{}, "01", status.UnAvailableServiceError, err
	}
	if repo.DBDS == nil {
		return enrollmentSchema.DeactivateEnrollmentResponse{}, "02", status.StatusBadRequest, err
	}
	cancel, err, massage := repo.db().CancelEnrollment(ctx, req.Body)
	if err != nil {
		return enrollmentSchema.DeactivateEnrollmentResponse{}, "03", status.StatusInternalServerError, err
	}
	return enrollmentSchema.DeactivateEnrollmentResponse{Enrollment: cancel, Result: massage}, "", status.StatusOK, nil
}

func (repo *Repository) ListEnrollment(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.ListEnrollmentsRequest]) (res enrollmentSchema.ListEnrollmentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return enrollmentSchema.ListEnrollmentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return enrollmentSchema.ListEnrollmentResponse{}, "02", status.StatusBadRequest, err
	}
	list, total, err := repo.db().ListEnrollment(ctx, req.Body)
	if err != nil {
		return enrollmentSchema.ListEnrollmentResponse{}, "03", status.StatusInternalServerError, err
	}
	return enrollmentSchema.ListEnrollmentResponse{Enrollments: list, TotalCount: total}, "", status.StatusOK, nil
}

func (repo *Repository) ListStudentCourse(ctx context.Context, req commonSchema.BaseRequest[enrollmentSchema.ListStudentCoursesRequest]) (res enrollmentSchema.ListStudentCoursesResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return enrollmentSchema.ListStudentCoursesResponse{}, "01", status.StatusInternalServerError, err
	}
	if repo.DBDS == nil {
		return enrollmentSchema.ListStudentCoursesResponse{}, "02", status.StatusBadRequest, err
	}
	list, err := repo.db().ListStudentCourses(ctx, req.Body)
	fmt.Println(list)
	if err != nil {
		return enrollmentSchema.ListStudentCoursesResponse{}, "03", status.StatusInternalServerError, err
	}
	return enrollmentSchema.ListStudentCoursesResponse{Enrollments: list}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.EnrollmentDS {
	return repo.DBDS
}
