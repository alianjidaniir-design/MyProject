package Registrations

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"MyProject/midddleware/authz"
	"MyProject/models/Registrations/dataSources"
	"MyProject/models/Registrations/dataSources/mysqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"sync"

	"github.com/gofiber/fiber/v2"
)

type Repository struct {
	DBDS     dataSources.RegistrationDS
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

	newEnr, err := mysqlDS.NewRegisterDBDS(dsn.RegistrationTableName, db)
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

func (repo *Repository) CreateRegistration(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.RegisterStudentRequest], c *fiber.Ctx) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.RegisterStudentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.RegisterStudentResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	role := authz.GetRoleName(c)
	studentID := authz.GetUserID(c)
	create, err := repo.db().RegistrationsStudent(ctx, req.Body, role, studentID)
	if err != nil {
		return registrationSchema.RegisterStudentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.RegisterStudentResponse{Information: create}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.DetailStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.DetailStudentResponse{}, "05", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.DetailStudentResponse{}, "06", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	get, err := repo.db().GetRegisterStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.DetailStudentResponse{}, "07", status.StatusInternalServerError, err
	}
	return registrationSchema.DetailStudentResponse{Information: get}, "", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.DetailStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.DetailStudentResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.DetailStudentResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	update, err := repo.db().UpdateRegisterStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.DetailStudentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.DetailStudentResponse{Information: update}, "", status.StatusOK, nil
}
func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest], c *fiber.Ctx) (res registrationSchema.DeleteStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.DeleteStudentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.DeleteStudentResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	role := authz.GetRoleName(c)
	studentID := authz.GetUserID(c)
	deleted, err := repo.db().DeleteRegisterStudent(ctx, req.Body, role, studentID)
	if err != nil {
		return registrationSchema.DeleteStudentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.DeleteStudentResponse{Information: deleted, Massage: "deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.SelectPageRegisteredStudentsRequest]) (res registrationSchema.ListStudentsResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ListStudentsResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.ListStudentsResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	list, total, err := repo.db().ListAllRegisterStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.ListStudentsResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ListStudentsResponse{List: list, Total: total}, "", status.StatusOK, nil
}
func (repo *Repository) Cancel(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.CancelRegistrationResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.CancelRegistrationResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.CancelRegistrationResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	cancel, err := repo.db().CancelRegisterStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.CancelRegistrationResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.CancelRegistrationResponse{Information: cancel}, "", status.StatusOK, nil
}

func (repo *Repository) ListStudents(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.ListStudentsRequest]) (res registrationSchema.ListStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ListStudentResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.ListStudentResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	student, tot, err := repo.db().ListStudentsOffering(ctx, req.Body)
	if err != nil {
		return registrationSchema.ListStudentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ListStudentResponse{List: student, Total: tot}, "", status.StatusOK, nil
}

func (repo *Repository) ListOfferings(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.ListOfferingRequest]) (res registrationSchema.ListOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ListOfferingResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.ListOfferingResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	offering, totalAll, err := repo.db().ListOfferingsStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.ListOfferingResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ListOfferingResponse{List: offering, Total: totalAll}, "", status.StatusOK, nil
}

func (repo *Repository) ListClassesStudent(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.Pages], c *fiber.Ctx) (res registrationSchema.ClassSchedule, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ClassSchedule{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.ClassSchedule{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	studentID := authz.GetUserID(c)
	class, tot, page, err := repo.db().ListClassesStudent(ctx, req.Body, studentID)
	if err != nil {
		return registrationSchema.ClassSchedule{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ClassSchedule{MyClasses: class, Total: tot, Page: page}, "", status.StatusOK, nil

}

func (repo *Repository) ListClassesTeacher(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.Pages], c *fiber.Ctx) (res registrationSchema.ClassesTeacher, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.ClassesTeacher{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.ClassesTeacher{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	teacherID := authz.GetUserID(c)
	detail, tot, page, err := repo.db().ListClassesTeacher(ctx, req.Body, teacherID)
	if err != nil {
		return registrationSchema.ClassesTeacher{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.ClassesTeacher{MyClasses: detail, Total: tot, Page: page}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.RegistrationDS {
	return repo.DBDS
}
