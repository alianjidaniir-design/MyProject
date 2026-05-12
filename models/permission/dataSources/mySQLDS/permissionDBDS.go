package mySQLDS

import (
	"MyProject/apiSchema/permissionSchema"
	"MyProject/models/permission/dataModel"
	"MyProject/models/permission/dataSources"
	"MyProject/pkg/pagination"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type PermissionDBDS struct {
	tableName string
	db        *sql.DB
}

func NewPermissionDBDS(tableName string, db *sql.DB) (dataSources.PermissionDS, error) {
	ff := &PermissionDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *PermissionDBDS) CreatePermission(ctx context.Context, req permissionSchema.CreatePermissionReq) (res dataModel.Permission, err error) {
	if req.Name == "" || len(req.Name) > 127 {
		return dataModel.Permission{}, errors.New("name length must be between 127 characters and name is not empty")
	}
	var tot int
	countName := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name = ?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countName, req.Name).Scan(&tot)
	if err != nil {
		return dataModel.Permission{}, err
	} else if tot > 1 {
		return dataModel.Permission{}, errors.New("there is already a role with this name")
	}
	insertQuery := fmt.Sprintf("INSERT INTO `%s` (`name`) VALUES (?)", ds.tableName)
	lastID, err := ds.db.ExecContext(ctx, insertQuery, req.Name)
	if err != nil {
		return dataModel.Permission{}, err
	}
	newID, err := lastID.LastInsertId()
	if err != nil {
		return dataModel.Permission{}, err
	}
	return ds.selectRoles(ctx, newID)

}

func (ds *PermissionDBDS) GetPermission(ctx context.Context, req permissionSchema.GetPermissionReq) (res dataModel.Permission, err error) {
	err = ds.checkPermission(ctx, req.ID)
	if err != nil {
		return dataModel.Permission{}, err
	}
	return ds.selectRoles(ctx, req.ID)
}

func (ds *PermissionDBDS) DeletePermission(ctx context.Context, req permissionSchema.GetPermissionReq) (res dataModel.Permission, err error) {
	err = ds.checkPermission(ctx, req.ID)
	if err != nil {
		return dataModel.Permission{}, err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE ID = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return dataModel.Permission{}, err
	}
	return dataModel.Permission{}, nil
}

func (ds *PermissionDBDS) ListPermissions(ctx context.Context, req permissionSchema.Pagination) (res []dataModel.Permission, total int, err error) {
	var permissions []dataModel.Permission
	page, pageSize, err := pagination.CheckPage(req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	limit := pageSize
	offset := (page - 1) * limit
	var tot int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&tot)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s ORDER BY ID LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission dataModel.Permission
		err = rows.Scan(&permission.ID, &permission.Name)
		if err != nil {
			return nil, 0, err
		}
		permissions = append(permissions, permission)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return permissions, tot, nil
}

func (ds *PermissionDBDS) selectRoles(ctx context.Context, ID int64) (dataModel.Permission, error) {
	var role dataModel.Permission
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=? ORDER BY id ", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&role.ID, &role.Name)
	if err != nil {
		return role, err
	}
	return role, nil
}

func (ds *PermissionDBDS) checkPermission(ctx context.Context, ID int64) error {
	var check bool
	checking := `
SELECT EXISTS (SELECT 1 FROM ` + ds.tableName + ` where id=? )`
	err := ds.db.QueryRowContext(ctx, checking, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("permission does not exist")
	}
	return nil
}
