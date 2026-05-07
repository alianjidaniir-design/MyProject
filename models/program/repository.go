package program

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/programSchema"
	"MyProject/models/program/dataSources"
	"MyProject/models/program/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.ProgramDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepo() {
	cfg, err := mySQLDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	open, err := mySQLDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	newProgram, err := mySQLDS.NewProgramDBDS(cfg.ProgramTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newProgram}
	log.Printf("successfully init repository")
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[programSchema.CreateProgramRequest]) (res programSchema.DetailProgramResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return programSchema.DetailProgramResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return programSchema.DetailProgramResponse{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateProgram(ctx, req.Body)
	if err != nil {
		return programSchema.DetailProgramResponse{}, "03", status.StatusBadRequest, err
	}
	return programSchema.DetailProgramResponse{Detail: create, Massage: "successfully created"}, "0", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.ProgramDS {
	return repo.DBDS
}
