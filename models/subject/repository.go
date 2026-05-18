package subject

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/subjectSchema"
	"MyProject/models/subject/dataSources"
	"MyProject/models/subject/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"sync"
)

type Repository struct {
	DBDS     dataSources.SubjectDS
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
	newSub, err := mySQLDS.NewSubjectDBDS(cfg.SubjectableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newSub, initRepo: err}

}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[subjectSchema.CreateSubject]) (res subjectSchema.DetailSubject, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return subjectSchema.DetailSubject{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return subjectSchema.DetailSubject{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateSubject(ctx, req.Body)
	if err != nil {
		return subjectSchema.DetailSubject{}, "03", status.StatusBadRequest, err
	}
	return subjectSchema.DetailSubject{Detail: create}, "", status.StatusOK, nil

}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[subjectSchema.GetSubject]) (res subjectSchema.DetailSubject, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return subjectSchema.DetailSubject{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return subjectSchema.DetailSubject{}, "02", status.StatusInternalServerError, err
	}
	get, err := repo.db().GetSubject(ctx, req.Body)
	if err != nil {
		return subjectSchema.DetailSubject{}, "03", status.StatusBadRequest, err
	}
	return subjectSchema.DetailSubject{Massage: "Detail Subject ", Detail: get}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[subjectSchema.GetSubject]) (res subjectSchema.DetailSubject, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return subjectSchema.DetailSubject{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return subjectSchema.DetailSubject{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeleteSubject(ctx, req.Body)
	if err != nil {
		return subjectSchema.DetailSubject{}, "03", status.StatusBadRequest, err
	}
	return subjectSchema.DetailSubject{Massage: "deleted successfully", Detail: deleted}, "", status.StatusOK, nil
}

func (repo *Repository)	List(ctx context.Context , req commonSchema.BaseRequest[subjectSchema.Pagination]) (res subjectSchema.ListSubjects, errStr string, code int, err error){
	if repo.initRepo != nil {
		return subjectSchema.ListSubjects{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return subjectSchema.ListSubjects{}, "02", status.StatusInternalServerError, err
	}
	list , tot , err:=repo.db().ListSubjects(ctx, req.Body)
	if err != nil {
		return subjectSchema.ListSubjects{}, "03", status.StatusBadRequest, err
	}
	return subjectSchema.ListSubjects{Massage: "List Subjects",List: list, Total: tot}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.SubjectDS {
	return repo.DBDS
}
