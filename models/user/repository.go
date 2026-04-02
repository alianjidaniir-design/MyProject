package user

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"MyProject/models/repositories"
	userDataSourses "MyProject/models/user/dataSourses"
	mysqlDataSource "MyProject/models/user/dataSourses/mySqlDS"
	"context"
	"log"
	"sync"
)

type Repository struct {
	dbDS    userDataSourses.UserDBDS
	initErr error
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[userSchema.ListRequest]) (userSchema.ListUser, string, int, error) {
	//TODO implement me
	panic("implement me")
}

var (
	once    sync.Once
	repoIns *Repository
)

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{}
		repoIns.initializeDataSources()
	})
	return repoIns
}

func init() {
	repositories.UserRepo = GetRepo()
}

func (repo *Repository) initializeDataSources() {
	mysqlDS, enabled, err := mysqlDataSource.NewUserDBDSFromEnv()
	if err != nil {
		repo.initErr = err
		return
	}
	if enabled {
		repo.dbDS = mysqlDS
		log.Printf("mysqlDataSource.NewUserDBDSFromEnv err:%v", mysqlDS)
	}
}

func (repo *Repository) db() userDataSourses.UserDBDS {
	return repo.dbDS
}
