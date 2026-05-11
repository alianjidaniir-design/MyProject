package book

import (
	"MyProject/apiSchema/bookSchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/book/dataSources"
	"MyProject/models/book/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"log"
	"net/http"
	"sync"
)

type Repository struct {
	DBDS     dataSources.BookDS
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
	newB, err := mySQLDS.NewBookDBDS(cfg.BookTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newB}
	log.Printf("done successfully")
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[bookSchema.RegistrationBook]) (res bookSchema.InformationBook, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return bookSchema.InformationBook{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return bookSchema.InformationBook{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().RegisterBook(ctx, req.Body)
	if err != nil {
		return bookSchema.InformationBook{}, "03", status.StatusInternalServerError, err
	}
	return bookSchema.InformationBook{Information: create}, "", http.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[bookSchema.GetCodeBook]) (res bookSchema.InformationBook, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return bookSchema.InformationBook{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return bookSchema.InformationBook{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeleteBook(ctx, req.Body)
	if err != nil {
		return bookSchema.InformationBook{}, "03", status.StatusInternalServerError, err
	}
	return bookSchema.InformationBook{Information: deleted, Massage: "deleted id successfully"}, "", http.StatusOK, nil
}
func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[bookSchema.GetCodeBook]) (res bookSchema.InformationBook, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return bookSchema.InformationBook{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return bookSchema.InformationBook{}, "02", status.StatusInternalServerError, err
	}
	get, err := repo.db().DetailBook(ctx, req.Body)
	if err != nil {
		return bookSchema.InformationBook{}, "03", status.StatusInternalServerError, err
	}
	return bookSchema.InformationBook{Information: get, Massage: "detail book"}, "", http.StatusOK, nil
}
func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[bookSchema.PaginationBook]) (res bookSchema.ListBooks, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return bookSchema.ListBooks{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return bookSchema.ListBooks{}, "02", status.StatusInternalServerError, err
	}
	list, tot, err := repo.db().ListBooks(ctx, req.Body)
	if err != nil {
		return bookSchema.ListBooks{}, "03", status.StatusInternalServerError, err
	}
	return bookSchema.ListBooks{Books: list, Total: tot}, "", http.StatusOK, nil
}

func (repo *Repository) db() dataSources.BookDS {
	return repo.DBDS
}
