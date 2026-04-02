package user

import (
	userDataSourses "MyProject/models/user/dataSourses"
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
	repoIns = &Repository{}
}
func GetRepo() *Repository {
	once.Do(initRepoIns)
	return repoIns
}

func (repo *Repository) db() userDataSourses.UserDBDS {
	return repo.dbDS
}
