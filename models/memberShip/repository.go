package memberShip

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/membershipSchema"
	"MyProject/models/memberShip/dataSources"
	"MyProject/models/memberShip/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"log"
	"sync"
)

type Repository struct {
	DBDS     dataSources.MembershipDS
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
	newMembership, err := mySQLDS.NewMembershipDBDS(cfg.MembershipTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newMembership}
	log.Printf("successfully init repository %+v", repo)
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.CreateMembershipRequest]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return membershipSchema.DetailMembershipSchema{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return membershipSchema.DetailMembershipSchema{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateMembership(ctx, req.Body)
	if err != nil {
		return membershipSchema.DetailMembershipSchema{}, "03", status.StatusBadRequest, err
	}
	return membershipSchema.DetailMembershipSchema{MemberShip: create, Massage: "created successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.MembershipDS {
	return repo.DBDS
}
