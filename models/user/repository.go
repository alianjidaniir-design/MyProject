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
	dbDS    userDataSourses.UserDBDS
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

	repoIns = &Repository{}
}

func GetRepoIns() *Repository {
	once.Do(initRepoIns)
	return repoIns
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.LoginRequest]) (res userSchema.ResponseUser, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.ResponseUser{}, "13", status.StatusUnauthorized, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.ResponseUser{}, "14", status.UnAvailableServiceError, errors.New("student datasourse not configured")
	}

	createdUser, err := repo.db().CreateStudent(ctx, req.Body)
	if err != nil {
		return userSchema.ResponseUser{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.ResponseUser{User: createdUser}, "", status.StatusOK, nil
}

func (repo *Repository) db() userDataSourses.UserDBDS {
	return repo.dbDS
}
