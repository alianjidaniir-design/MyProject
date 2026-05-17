package publisher

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/publisherSchema"
	"MyProject/models/publisher/dataSources"
	"MyProject/models/publisher/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"net/http"
	"sync"
)

type Repository struct {
	DBDS     dataSources.PublisherDS
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
	newPub, err := mySQLDS.NewPublisherDBDS(cfg.PublisherTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newPub}
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[publisherSchema.CreatePublisher]) (res publisherSchema.DetailPublisher, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return publisherSchema.DetailPublisher{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return publisherSchema.DetailPublisher{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreatePublisher(ctx, req.Body)
	if err != nil {
		return publisherSchema.DetailPublisher{}, "03", status.StatusBadRequest, err
	}
	return publisherSchema.DetailPublisher{Detail: create}, "", http.StatusCreated, nil
}

func (repo *Repository) Detail(ctx context.Context, req commonSchema.BaseRequest[publisherSchema.GetPublisher]) (res publisherSchema.DetailPublisher, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return publisherSchema.DetailPublisher{}, "01", status.StatusBadRequest, repo.initRepo
	}
	if repo.DBDS == nil {
		return publisherSchema.DetailPublisher{}, "02", status.StatusInternalServerError, err
	}
	detail, err := repo.db().DetailPublisher(ctx, req.Body)
	if err != nil {
		return publisherSchema.DetailPublisher{}, "03", status.StatusBadRequest, err
	}
	return publisherSchema.DetailPublisher{Massage: "Detail Publisher", Detail: detail}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.PublisherDS {
	return repo.DBDS
}
