package author

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/author/dataSources"
	"MyProject/models/author/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.AuthorDS
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
	newAt, err := mySQLDS.NewAuthorDBDS(cfg.AuthorTableName, open)

	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newAt}
	log.Printf("successfully initialized author table")
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) CreateAuthor(ctx context.Context, req commonSchema.BaseRequest[authorSchema.CreateAuthor]) (res authorSchema.DetailAuthor, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return authorSchema.DetailAuthor{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return authorSchema.DetailAuthor{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateAuthor(ctx, req.Body)
	if err != nil {
		return authorSchema.DetailAuthor{}, "03", status.StatusBadRequest, err
	}
	return authorSchema.DetailAuthor{Detail: create, Massage: "created successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.AuthorDS {
	return repo.DBDS
}
