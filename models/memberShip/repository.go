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

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.GetIDMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return membershipSchema.DetailMembershipSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return membershipSchema.DetailMembershipSchema{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeleteMembership(ctx, req.Body)
	if err != nil {
		return membershipSchema.DetailMembershipSchema{}, "03", status.StatusBadRequest, err
	}
	return membershipSchema.DetailMembershipSchema{MemberShip: deleted}, "", status.StatusOK, nil
}

func (repo *Repository) Update(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.UpdateMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return membershipSchema.DetailMembershipSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return membershipSchema.DetailMembershipSchema{}, "02", status.StatusInternalServerError, err
	}
	update, err := repo.db().UpdateMembership(ctx, req.Body)
	if err != nil {
		return membershipSchema.DetailMembershipSchema{}, "03", status.StatusBadRequest, err
	}
	return membershipSchema.DetailMembershipSchema{MemberShip: update, Massage: "updated successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) DeActive(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.GetIDMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return membershipSchema.DetailMembershipSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return membershipSchema.DetailMembershipSchema{}, "02", status.StatusInternalServerError, err
	}
	deActive, err := repo.db().DeActiveMembership(ctx, req.Body)
	if err != nil {
		return membershipSchema.DetailMembershipSchema{}, "03", status.StatusBadRequest, err
	}
	return membershipSchema.DetailMembershipSchema{MemberShip: deActive, Massage: "deactivated successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.GetIDMembership]) (res membershipSchema.DetailMembershipSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return membershipSchema.DetailMembershipSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return membershipSchema.DetailMembershipSchema{}, "02", status.StatusInternalServerError, err
	}
	detail, err := repo.db().DetailMembership(ctx, req.Body)
	if err != nil {
		return membershipSchema.DetailMembershipSchema{}, "03", status.StatusBadRequest, err
	}
	return membershipSchema.DetailMembershipSchema{MemberShip: detail, Massage: "detail membership student"}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[membershipSchema.PaginationMemberShip]) (res membershipSchema.ListMembershipSchema, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return membershipSchema.ListMembershipSchema{}, "01", status.StatusInternalServerError, repo.initRepo
	}
	if repo.DBDS == nil {
		return membershipSchema.ListMembershipSchema{}, "02", status.StatusInternalServerError, err
	}
	list, totalRows, err := repo.db().ListMembership(ctx, req.Body)
	if err != nil {
		return membershipSchema.ListMembershipSchema{}, "03", status.StatusBadRequest, err
	}
	return membershipSchema.ListMembershipSchema{List: list, Total: totalRows}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.MembershipDS {
	return repo.DBDS
}
