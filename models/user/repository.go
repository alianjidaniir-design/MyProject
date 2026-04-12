package user

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	userDataSourses "MyProject/models/user/dataSourses"
	"MyProject/models/user/dataSourses/mySqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

type Repository struct {
	dbDS    userDataSourses.UserDB
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

	userDNInstance, err := mySqlDS.NewUsersDBDS(DBConn, cfg.StudentTableName)
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

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.LoginRequest]) (res userSchema.ResponseUser, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.ResponseUser{}, "13", status.StatusUnauthorized, repo.initErr
	}
	if repo.dbDS == nil {
		return userSchema.ResponseUser{}, "14", status.UnAvailableServiceError, errors.New("student dataSource not configured")
	}

	createdUser, err := repo.db().CreateStudent(ctx, req.Body)
	if err != nil {
		return userSchema.ResponseUser{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.ResponseUser{User: createdUser}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[userSchema.ListRequest]) (res userSchema.ListUser, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.ListUser{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.ListUser{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	lists, total, err := repo.db().ReadStudent(ctx, req.Body)
	if err != nil {
		return userSchema.ListUser{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.ListUser{Users: lists, Total: total}, "", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[userSchema.UpdateUserRequest]) (res userSchema.UpdateResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.UpdateResponse{}, "14", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.UpdateResponse{}, "15", status.StatusInternalServerError, errors.New("bad")
	}
	updatedUser, err := repo.db().UpdateStudent(ctx, req.Body)
	if err != nil {
		return userSchema.UpdateResponse{}, "16", status.UnAvailableServiceError, err
	}
	return userSchema.UpdateResponse{User: updatedUser}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[userSchema.GetRequest]) (res userSchema.GetResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.GetResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.GetResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}

	getIng, err := repo.db().GetStudent(ctx, req.Body)
	if err != nil {
		return userSchema.GetResponse{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.GetResponse{User: getIng}, "", status.StatusOK, nil
}

func (repo *Repository) SoftDelete(ctx context.Context, req commonSchema.BaseRequest[userSchema.SoftDeleteRequest]) (res userSchema.SoftDeleteResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.SoftDeleteResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.SoftDeleteResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	soft, err := repo.db().SoftDeleteStudent(ctx, req.Body)
	if err != nil {
		return userSchema.SoftDeleteResponse{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.SoftDeleteResponse{User: soft}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[userSchema.DeleteRequest]) (res userSchema.DeleteResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.DeleteResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.DeleteResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	deletedUser, err := repo.db().DeleteStudent(ctx, req.Body)
	if err != nil {
		return userSchema.DeleteResponse{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.DeleteResponse{User: deletedUser}, "", status.StatusOK, nil
}

func (repo *Repository) db() userDataSourses.UserDB {
	return repo.dbDS
}
