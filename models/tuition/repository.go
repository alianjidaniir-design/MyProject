package tuition

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/models/tuition/dataSources"
	"MyProject/models/tuition/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"fmt"
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
		repo = &Repository{initRepo: errors.New("problem loading config")}
		return
	}
	open, err := mySQLDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	newTui, err := mySQLDS.NewTuitionDBDS(cfg.TuitionTableName, open)
	if err != nil {
		repo = &Repository{initRepo: errors.New(err.Error())}
		return
	}
	repo = &Repository{DBDS: newTui}
	fmt.Println(newTui)
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

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[tuitionSchema.UpdateTuition]) (res tuitionSchema.MassageUpdateTuition, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return tuitionSchema.MassageUpdateTuition{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return tuitionSchema.MassageUpdateTuition{}, "02", status.StatusInternalServerError, err
	}
	update, err := repo.db().UpdateTuition(ctx, req.Body)
	if err != nil {
		return tuitionSchema.MassageUpdateTuition{}, "03", status.StatusBadRequest, err
	}
	return tuitionSchema.MassageUpdateTuition{Detail: update, Massage: "updated successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.TuitionDS {
	return repo.DBDS
}
