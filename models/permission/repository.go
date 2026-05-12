package permission

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/permissionSchema"
	"MyProject/models/permission/dataSources"
	"MyProject/models/permission/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"net/http"
	"sync"
)

type Repository struct {
	DBDS     dataSources.PermissionDS
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
	newPer, err := mySQLDS.NewPermissionDBDS(cfg.PermissionTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
	}
	repo = &Repository{DBDS: newPer}
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.CreatePermissionReq]) (res permissionSchema.DetailPermission, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return permissionSchema.DetailPermission{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return permissionSchema.DetailPermission{}, "02", status.StatusInternalServerError, err
	}
	create, err := repo.db().CreatePermission(ctx, req.Body)
	if err != nil {
		return permissionSchema.DetailPermission{}, "03", status.StatusBadRequest, err
	}
	return permissionSchema.DetailPermission{Permission: create}, "", http.StatusOK, nil
}

func (repo *Repository) Get(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.GetPermissionReq]) (res permissionSchema.DetailPermission, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return permissionSchema.DetailPermission{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return permissionSchema.DetailPermission{}, "02", status.StatusInternalServerError, err
	}
	get, err := repo.db().GetPermission(ctx, req.Body)
	if err != nil {
		return permissionSchema.DetailPermission{}, "03", status.StatusBadRequest, err
	}
	return permissionSchema.DetailPermission{Permission: get, Massage: "Detail Permission"}, "", http.StatusOK, nil
}

func (repo *Repository) Delete(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.GetPermissionReq]) (res permissionSchema.DetailPermission, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return permissionSchema.DetailPermission{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return permissionSchema.DetailPermission{}, "02", status.StatusInternalServerError, err
	}
	deleted, err := repo.db().DeletePermission(ctx, req.Body)
	if err != nil {
		return permissionSchema.DetailPermission{}, "03", status.StatusBadRequest, err
	}
	return permissionSchema.DetailPermission{Permission: deleted, Massage: "deleted permission successfully"}, "", http.StatusOK, nil
}

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[permissionSchema.Pagination]) (res permissionSchema.ListPermission, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return permissionSchema.ListPermission{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return permissionSchema.ListPermission{}, "02", status.StatusInternalServerError, err
	}
	list, tot, err := repo.db().ListPermissions(ctx, req.Body)
	if err != nil {
		return permissionSchema.ListPermission{}, "03", status.StatusBadRequest, err
	}
	return permissionSchema.ListPermission{Massage: "List Permission", Permission: list, TotalPermission: tot}, "", http.StatusOK, nil
}

func (repo *Repository) db() dataSources.PermissionDS {
	return repo.DBDS
}
