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
	dbconn, err := mySqlDS.Open(cfg)
	if err != nil {
		repoIns = &Repository{initErr: fmt.Errorf("failed to open config: %v", err)}
		log.Printf("Error opening DB connection: %v", err)
		return
	}

	userDNInstance, err := mySqlDS.NewUsersDBDS(dbconn, cfg.StudentTableName)
	if err != nil {
		_ = dbconn.Close()
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
		return userSchema.ResponseUser{}, "14", status.UnAvailableServiceError, errors.New("student datasourse not configured")
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
	listus, total, err := repo.db().ReadStudent(ctx, req.Body)
	if err != nil {
		return userSchema.ListUser{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.ListUser{Users: listus, Total: total}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[userSchema.GetRequest]) (res userSchema.GetResponse, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.GetResponse{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.GetResponse{}, "11", status.StatusInternalServerError, errors.New("bad")
	}

	geting, err := repo.db().GetStudent(ctx, req.Body)
	if err != nil {
		return userSchema.GetResponse{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.GetResponse{User: geting}, "", status.StatusOK, nil
}

func (repo *Repository) db() userDataSourses.UserDB {
	return repo.dbDS
}
