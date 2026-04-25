package profile

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/profileSchema"
	"MyProject/models/profile/dataSources"
	"MyProject/models/profile/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"errors"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.ProfileDS
	initRepo error
}

var (
	once sync.Once
	repo *Repository
)

func initRepository() {
	cfg, err := mySQLDS.LoadConfig()
	if err != nil {
		repo = &Repository{initRepo: errors.New("failed to load config")}
		return
	}
	open, err := mySQLDS.Open(cfg)
	if err != nil {
		repo = &Repository{initRepo: errors.New("failed to open database")}
	}
	newProfile, err := mySQLDS.NewProfileDBDS(cfg.TableName, open)
	if err != nil {
		repo = &Repository{initRepo: errors.New("failed to initialize database")}
	}
	repo = &Repository{DBDS: newProfile}
	log.Println("initialized database successfully")
}

func GetRepo() *Repository {
	once.Do(initRepository)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[profileSchema.CreateScoresReq]) (res profileSchema.InformationResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return profileSchema.InformationResponse{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return profileSchema.InformationResponse{}, "02", status.StatusUnauthorized, errors.New("db is nil")
	}
	created, err := repo.db().CreateScoreStudent(ctx, req.Body)
	if err != nil {
		return profileSchema.InformationResponse{}, "03", status.StatusBadRequest, err
	}
	return profileSchema.InformationResponse{Details: created, Message: "created successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[profileSchema.ListAllScoresReq]) (res profileSchema.ListAllScoresResp, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return profileSchema.ListAllScoresResp{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return profileSchema.ListAllScoresResp{}, "02", status.StatusUnauthorized, errors.New("db is nil")
	}
	list, totalPage, err := repo.db().ListScoresStudents(ctx, req.Body)
	if err != nil {
		return profileSchema.ListAllScoresResp{}, "03", status.StatusBadRequest, err
	}
	return profileSchema.ListAllScoresResp{Students: list, Total: totalPage}, "", status.StatusOK, nil
}

func (repo *Repository) Summery(ctx context.Context, req commonSchema.BaseRequest[profileSchema.ListAllScoresReq]) (res profileSchema.StudentsSummeryResponse, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return profileSchema.StudentsSummeryResponse{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return profileSchema.StudentsSummeryResponse{}, "02", status.StatusUnauthorized, errors.New("db is nil")
	}
	list, totalPages, err := repo.db().ListSummeryStudents(ctx, req.Body)
	if err != nil {
		return profileSchema.StudentsSummeryResponse{}, "03", status.StatusBadRequest, err
	}
	return profileSchema.StudentsSummeryResponse{Summery: list, Total: totalPages}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[profileSchema.GetScoresReq]) (res profileSchema.DetailProfileStudent, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return profileSchema.DetailProfileStudent{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return profileSchema.DetailProfileStudent{}, "02", status.StatusUnauthorized, errors.New("db is nil")
	}
	get, err := repo.db().GetStudent(ctx, req.Body)
	if err != nil {
		return profileSchema.DetailProfileStudent{}, "03", status.StatusBadRequest, err
	}
	return profileSchema.DetailProfileStudent{Detail: get}, "", status.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[profileSchema.DeleteScoresReq]) (res profileSchema.DeleteProfileScoresResp, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return profileSchema.DeleteProfileScoresResp{}, "01", status.StatusUnauthorized, repo.initRepo
	}
	if repo.DBDS == nil {
		return profileSchema.DeleteProfileScoresResp{}, "02", status.StatusUnauthorized, errors.New("db is nil")
	}
	err = repo.db().DeleteProfile(ctx, req.Body)
	if err != nil {
		return profileSchema.DeleteProfileScoresResp{}, "03", status.StatusBadRequest, err
	}
	return profileSchema.DeleteProfileScoresResp{Message: "deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.ProfileDS {
	return repo.DBDS
}
