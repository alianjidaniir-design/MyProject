package enrollment

import (
	"MyProject/models/enrollment/dataSources"
	"MyProject/models/enrollment/dataSources/mysqlDS"
	"MyProject/models/user/dataSourses/mySqlDS"
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
	dsn, err := mySqlDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	db, err := mySqlDS.Open(dsn)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	defer db.Close()

	insta, err := mysqlDS.NewEnrollmentDBDS(dsn.StudentTableName, db)
	if err != nil {
		repo = &Repository{initRepo: err}
	}
	repo = &Repository{DBDS: insta}

}

func GetRepo() *Repository {
	onceEnrollment.Do(instance)
	return repo
}
