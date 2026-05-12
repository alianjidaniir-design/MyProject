package role

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/roleSchema"
	"MyProject/models/role/dataSources"
	"MyProject/models/role/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"sync"
)

type Repository struct {
	DBDS     dataSources.RoleDS
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
	newRole, err := mySQLDS.NewRoleDBDS(cfg.RoleTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newRole}
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[roleSchema.CreateRoleRequest]) (res roleSchema.DetailRole, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return roleSchema.DetailRole{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return roleSchema.DetailRole{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreateRole(ctx, req.Body)
	if err != nil {
		return roleSchema.DetailRole{}, "03", status.StatusBadRequest, err
	}
	return roleSchema.DetailRole{Role: create}, "", status.StatusOK, nil

}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[roleSchema.GetRoleRequest]) (res roleSchema.DetailRole, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return roleSchema.DetailRole{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return roleSchema.DetailRole{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeleteRole(ctx, req.Body)
	if err != nil {
		return roleSchema.DetailRole{}, "03", status.StatusBadRequest, err
	}
	return roleSchema.DetailRole{Role: deleted, Massage: "role deleted successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[roleSchema.GetRoleRequest]) (res roleSchema.DetailRole, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return roleSchema.DetailRole{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return roleSchema.DetailRole{}, "02", status.StatusInternalServerError, err
	}
	get, err := repo.db().GetRole(ctx, req.Body)
	if err != nil {
		return roleSchema.DetailRole{}, "03", status.StatusBadRequest, err
	}
	return roleSchema.DetailRole{Role: get, Massage: "Detail Roles"}, "", status.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[roleSchema.Pagination]) (res roleSchema.ListRole, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return roleSchema.ListRole{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return roleSchema.ListRole{}, "02", status.StatusInternalServerError, err
	}
	list, tot, err := repo.db().ListRoles(ctx, req.Body)
	if err != nil {
		return roleSchema.ListRole{}, "03", status.StatusBadRequest, err
	}
	return roleSchema.ListRole{DetailRole: list, Total: tot}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.RoleDS {
	return repo.DBDS
}
