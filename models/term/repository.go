package term

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/termSchema"
	"MyProject/models/term/dataSources"
	"MyProject/models/term/dataSources/mySqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"sync"
)

type Repository struct {
	DBDS     dataSources.TermDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepository() {
	cfg, err := mySqlDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	load, err := mySqlDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: errors.New("problem opening database")}
		return
	}
	newterm, err := mySqlDS.NewTermDBDS(cfg.TermTableName, load)
	if err != nil {
		repo = &Repository{initRepo: err}
	}
	repo = &Repository{DBDS: newterm}
}

func GetRepo() *Repository {
	once.Do(initRepository)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[termSchema.CreateTerm]) (res termSchema.InformationTerm, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return termSchema.InformationTerm{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return termSchema.InformationTerm{}, "02", status.StatusBadRequest, err
	}
	created, err := repo.db().CreateTerm(ctx, req.Body)
	if err != nil {
		return termSchema.InformationTerm{}, "03", status.UnAvailableServiceError, err
	}
	return termSchema.InformationTerm{Term: created}, "", status.StatusOK, nil

}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[termSchema.ListTerm]) (res termSchema.ListTerms, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return termSchema.ListTerms{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return termSchema.ListTerms{}, "02", status.StatusBadRequest, err
	}
	list, totalPages, err := repo.db().ListTerms(ctx, req.Body)
	if err != nil {
		return termSchema.ListTerms{}, "03", status.UnAvailableServiceError, err
	}
	return termSchema.ListTerms{Term: list, Total: totalPages}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[termSchema.DeleteTerm]) (res termSchema.DeletedTerm, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return termSchema.DeletedTerm{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return termSchema.DeletedTerm{}, "02", status.StatusBadRequest, err
	}
	_, err = repo.db().DeleteTerms(ctx, req.Body)
	if err != nil {
		return termSchema.DeletedTerm{}, "03", status.UnAvailableServiceError, err
	}
	return termSchema.DeletedTerm{Message: "user deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.TermDS {
	return repo.DBDS
}
