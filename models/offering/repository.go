package offering

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/offeringSchema"
	"MyProject/models/offering/dataSources"
	"MyProject/models/offering/dataSources/mySqlDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"sync"
)

type Repository struct {
	DBDS     dataSources.OfferingDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepository() {
	cfg, err := mySqlDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: errors.New("you can not load config" + err.Error())}
		return
	}
	open, err := mySqlDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: errors.New("you can not open config" + err.Error())}
		return
	}
	newOffer, err := mySqlDS.NewOfferingDBDS(cfg.TableName, open)
	if err != nil {
		repo = &Repository{initRepo: errors.New(err.Error())}
		return
	}
	repo = &Repository{DBDS: newOffer}
}

func GetRepository() *Repository {
	once.Do(initRepository)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.CreateOfferingRequest]) (res offeringSchema.CreateOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.CreateOfferingResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.CreateOfferingResponse{}, "02", status.StatusBadRequest, err
	}
	create, err := repo.db().CreateOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.CreateOfferingResponse{}, "03", status.StatusInternalServerError, err
	}
	return offeringSchema.CreateOfferingResponse{Specification: create}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.ListOfferingsRequest]) (res offeringSchema.ListOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.ListOfferingResponse{}, "04", status.StatusBadRequest, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.ListOfferingResponse{}, "05", status.StatusBadRequest, err
	}
	list, total, err := repo.db().ListOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.ListOfferingResponse{}, "06", status.StatusInternalServerError, err
	}
	return offeringSchema.ListOfferingResponse{Offerings: list, TotalCount: total}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[offeringSchema.GetRowOfferingRequest]) (res offeringSchema.DetailOfferingResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return offeringSchema.DetailOfferingResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return offeringSchema.DetailOfferingResponse{}, "02", status.StatusBadRequest, err
	}
	get, err := repo.db().GetOffering(ctx, req.Body)
	if err != nil {
		return offeringSchema.DetailOfferingResponse{}, "03", status.StatusInternalServerError, err
	}
	return offeringSchema.DetailOfferingResponse{Specification: get}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.OfferingDS {
	return repo.DBDS
}
