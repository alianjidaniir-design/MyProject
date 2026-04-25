package tuition

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition/dataSources"
	"MyProject/models/tuition/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"sync"
)

type Repository struct {
	DBDS     dataSources.TuitionDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepository() {
	cfg, err := mySQLDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: errors.New(err.Error())}
	}
	open, err := mySQLDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: errors.New("there isa problem in opening")}
	}
	newTui, err := mySQLDS.NewTuitionDBDS(open, cfg.TermTableName)
	if err != nil {
		repo = &Repository{initRepo: errors.New(err.Error())}
	}
	repo = &Repository{DBDS: newTui}
}

func GetRepo() *Repository {
	once.Do(initRepository)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.CreateTuition]) (res tuitionSchema.InformationTuitionSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return tuitionSchema.InformationTuitionSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return tuitionSchema.InformationTuitionSchema{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateTuition(ctx, req.Body)
	if err != nil {
		return tuitionSchema.InformationTuitionSchema{}, "03", status.StatusBadRequest, err
	}
	return tuitionSchema.InformationTuitionSchema{Detail: create}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.TuitionDS {
	return repo.DBDS
}
