package user

import (
	"MyProject/models/repositories"
	userDataSourses "MyProject/models/user/dataSourses"
	"MyProject/models/user/dataSourses/memoryDS"
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

func GetRepo() *Repository {
	once.Do(func() {
		repoIns = &Repository{
			dbDS: memoryDS.NewTaskDBDS(1),
		}
	})
	return repoIns
}

func init() {
	repositories.UserRepo = GetRepo()
}

func (repo *Repository) db() userDataSourses.UserDBDS {
	return repo.dbDS
}
