package student

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/studentSchema"
	userDataSourses "MyProject/models/student/dataSourses"
	"MyProject/models/student/dataSourses/mySqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Repository struct {
	dbDS    userDataSourses.StudentDB
	initErr error
}

var (
	once    sync.Once
	repoIns *Repository
)

func initRepoIns() {

	cfg, err := mySqlDS.LoadConfig()
	if err != nil {
		repoIns = &Repository{initErr: fmt.Errorf("failed to load config: %v", err)}
		return
	}
	DBConn, err := mySqlDS.Open(cfg)
	if err != nil {
		repoIns = &Repository{initErr: fmt.Errorf("failed to open config: %v", err)}
		log.Printf("Error opening DB connection: %v", err)
		return
	}

	userDNInstance, err := mySqlDS.NewStudentDBDS(DBConn, cfg.StudentTableName)
	if err != nil {
		_ = DBConn.Close()
		repoIns = &Repository{initErr: fmt.Errorf("failed to connect to DB: %v", err)}
		log.Printf("Error opening DB connection: %v", err)
		return
	}

	repoIns = &Repository{dbDS: userDNInstance}
	log.Println("repository init success")
}

func GetRepoIns() *Repository {
	once.Do(initRepoIns)
	return repoIns
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[studentSchema.SignUpStudent]) (res studentSchema.UserResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.UserResponse{}, "13", status.StatusUnauthorized, repo.initErr
	}
	if repo.dbDS == nil {
		return studentSchema.UserResponse{}, "14", status.UnAvailableServiceError, errors.New("student dataSource not configured")
	}

	createdUser, err := repo.db().CreateStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.UserResponse{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.UserResponse{User: createdUser}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[studentSchema.ListRequest]) (res studentSchema.ListUser, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.ListUser{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.ListUser{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	lists, total, err := repo.db().ReadStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.ListUser{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.ListUser{Users: lists, Total: total}, "", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[studentSchema.UpdateUserRequest]) (res studentSchema.UpdateResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.UpdateResponse{}, "14", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.UpdateResponse{}, "15", status.StatusInternalServerError, errors.New("bad")
	}
	updatedUser, err := repo.db().UpdateStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.UpdateResponse{}, "16", status.UnAvailableServiceError, err
	}
	return studentSchema.UpdateResponse{User: updatedUser}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[studentSchema.GetRequest]) (res studentSchema.GetResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.GetResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.GetResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}

	getIng, err := repo.db().GetStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.GetResponse{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.GetResponse{User: getIng}, "", status.StatusOK, nil
}

func (repo *Repository) SoftDelete(ctx context.Context, req commonSchema.BaseRequest[studentSchema.SoftDeleteRequest]) (res studentSchema.SoftDeleteResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.SoftDeleteResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.SoftDeleteResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	soft, err := repo.db().SoftDeleteStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.SoftDeleteResponse{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.SoftDeleteResponse{User: soft}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[studentSchema.DeleteRequest]) (res studentSchema.DeleteResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.DeleteResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.DeleteResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	deletedUser, err := repo.db().DeleteStudent(ctx, req.Body)
	if err != nil {
		return studentSchema.DeleteResponse{}, "04", status.UnAvailableServiceError, err
	}
	return studentSchema.DeleteResponse{User: deletedUser}, "", status.StatusOK, nil
}

func (repo *Repository) Entry(ctx context.Context, req commonSchema.BaseRequest[studentSchema.LoginStudent]) (res studentSchema.StudentEntry, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.StudentEntry{}, "01", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.StudentEntry{}, "02", status.StatusInternalServerError, errors.New("bad")
	}
	access, refresh, err := repo.db().StudentEntry(ctx, req.Body)
	if err != nil {
		return studentSchema.StudentEntry{}, "03", status.UnAvailableServiceError, err
	}
	return studentSchema.StudentEntry{Massage: "login successfully", AccessToken: access, RefreshToken: refresh}, "", status.StatusOK, nil

}

func (repo *Repository) RefreshToken(ctx context.Context, req commonSchema.BaseRequest[studentSchema.RefreshTokenRequest]) (res studentSchema.RefreshTokenResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return studentSchema.RefreshTokenResponse{}, "01", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return studentSchema.RefreshTokenResponse{}, "02", status.StatusInternalServerError, errors.New("bad")
	}
	refreshToken, accessToken, err := repo.db().RefreshToken(ctx, req.Body)
	if err != nil {
		return studentSchema.RefreshTokenResponse{}, "03", status.UnAvailableServiceError, err
	}
	return studentSchema.RefreshTokenResponse{RefreshToken: refreshToken, AccessToken: accessToken}, "", status.StatusOK, nil
}

func (repo *Repository) db() userDataSourses.StudentDB {
	return repo.dbDS
}
