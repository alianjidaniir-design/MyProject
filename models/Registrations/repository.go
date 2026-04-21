package Registrations

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"MyProject/models/Registrations/dataSources"
	"MyProject/models/Registrations/dataSources/mysqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"sync"
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

	newEnr, err := mysqlDS.NewEnrollmentDBDS(dsn.RegistrationTableName, db)
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

func (repo *Repository) CreateRegistration(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.RegisterStudentRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.RegisterStudentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.RegisterStudentResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	create, err := repo.db().RegistrationsStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.RegisterStudentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.RegisterStudentResponse{Information: create}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.RegisterStudentResponse{}, "05", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.RegisterStudentResponse{}, "06", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	get, err := repo.db().GetRegisterStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.RegisterStudentResponse{}, "07", status.StatusInternalServerError, err
	}
	return registrationSchema.RegisterStudentResponse{Information: get}, "", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.RegisterStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.RegisterStudentResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.RegisterStudentResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	update, err := repo.db().UpdateRegisterStudent(ctx, req.Body)
	if err != nil {
		return registrationSchema.RegisterStudentResponse{}, "03", status.StatusInternalServerError, err
	}
	return registrationSchema.RegisterStudentResponse{Information: update}, "", status.StatusOK, nil
}
func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[registrationSchema.GetRegisteredStudentsRequest]) (res registrationSchema.DeleteStudentResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return registrationSchema.DeleteStudentResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return registrationSchema.DeleteStudentResponse{}, "02", status.StatusBadRequest, errors.New("DB DS not initialized")
	}
	deleted, err := repo.db().DeleteRegisterStudent(ctx, req.Body)
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

func (repo *Repository) db() dataSources.RegistrationDS {
	return repo.DBDS
}
