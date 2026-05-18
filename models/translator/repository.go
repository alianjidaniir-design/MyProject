package translator

import (
	"MyProject/apiSchema/commonSchema"
	translateSchema "MyProject/apiSchema/translatorSchema"
	"MyProject/models/translator/dataSources"
	"MyProject/models/translator/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.TranslatorDS
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
	newAt, err := mySQLDS.NewTranslatorDBDS(cfg.TranslatorTableName, open)

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

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[translateSchema.CreateTranslator]) (res translateSchema.DetailTranslator, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return translateSchema.DetailTranslator{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return translateSchema.DetailTranslator{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateTranslator(ctx, req.Body)
	if err != nil {
		return translateSchema.DetailTranslator{}, "03", status.StatusBadRequest, err
	}
	return translateSchema.DetailTranslator{Detail: create, Massage: "created successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[translateSchema.GetTranslator]) (res translateSchema.DetailTranslator, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return translateSchema.DetailTranslator{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return translateSchema.DetailTranslator{}, "02", status.StatusInternalServerError, err
	}
	get, err := repo.db().GetTranslatorAuthor(ctx, req.Body)
	if err != nil {
		return translateSchema.DetailTranslator{}, "03", status.StatusBadRequest, err
	}
	return translateSchema.DetailTranslator{Detail: get, Massage: "Detail Author"}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[translateSchema.GetTranslator]) (res translateSchema.DetailTranslator, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return translateSchema.DetailTranslator{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return translateSchema.DetailTranslator{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeleteTranslator(ctx, req.Body)
	if err != nil {
		return translateSchema.DetailTranslator{}, "03", status.StatusBadRequest, err
	}
	return translateSchema.DetailTranslator{Detail: deleted, Massage: "deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[translateSchema.Pagination]) (res translateSchema.ListTranslator, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return translateSchema.ListTranslator{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return translateSchema.ListTranslator{}, "02", status.StatusInternalServerError, err
	}
	list, tot, err := repo.db().ListTranslator(ctx, req.Body)
	if err != nil {
		return translateSchema.ListTranslator{}, "03", status.StatusBadRequest, err
	}
	return translateSchema.ListTranslator{Massage: "List author", List: list, Total: tot}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.TranslatorDS {
	return repo.DBDS
}
