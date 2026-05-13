package rolePermission

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/rolePermissionSchema"
	"MyProject/models/rolePermission/dataSources"
	"MyProject/models/rolePermission/dataSources/mySQLDS"
	"MyProject/statics/constants/status"
	"context"
	"sync"
)

type Repository struct {
	DBDS     dataSources.RolePermissionDS
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
	newRP, err := mySQLDS.NewRolePermissionDBDS(cfg.RolePermissionTableName, open)
	if err != nil {
		repo = &Repository{initRepo: err}
		return
	}
	repo = &Repository{DBDS: newRP}
}

func GetRepo() *Repository {
	once.Do(initRepo)
	return repo
}

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[rolePermissionSchema.CreateRolePermissionReq]) (res rolePermissionSchema.DetailRolePermission, errStr string, code int, err error) {
	if repo.initRepo != nil {
		return rolePermissionSchema.DetailRolePermission{}, "01", status.UnAvailableServiceError, repo.initRepo
	}
	if repo.DBDS == nil {
		return rolePermissionSchema.DetailRolePermission{}, "02", status.StatusInternalServerError, err
	}
	err = repo.db().CreatePermission(ctx, req.Body)
	if err != nil {
		return rolePermissionSchema.DetailRolePermission{}, "03", status.StatusBadRequest, err
	}
	return rolePermissionSchema.DetailRolePermission{Massage: "created successfully"}, "", status.StatusOK, nil
}

func (repo *Repository) db() dataSources.RolePermissionDS {
	return repo.DBDS
}
