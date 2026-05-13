package mySQLDS

import (
	"MyProject/apiSchema/rolePermissionSchema"
	"MyProject/models/rolePermission/dataModel"
	"MyProject/models/rolePermission/dataSources"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RolePermissionDBDS struct {
	tableName string
	db        *sql.DB
}

func NewRolePermissionDBDS(tableName string, db *sql.DB) (dataSources.RolePermissionDS, error) {
	ff := &RolePermissionDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *RolePermissionDBDS) CreatePermission(ctx context.Context, req rolePermissionSchema.CreateRolePermissionReq) (err error) {
	var checkRole, checkPermission bool
	checkingQuery := `
SELECT
 EXISTS (SELECT 1 FROM roles WHERE ID = ? ),
 EXISTS (SELECT 1 FROM permissions WHERE ID = ? )`
	err = ds.db.QueryRowContext(ctx, checkingQuery, req.RoleID, req.PermissionID).Scan(&checkRole, &checkPermission)
	if err != nil {
		return err
	}
	if checkRole == false || checkPermission == false {
		return errors.New("Role or Permission does not exist")
	}
	insertQuery := fmt.Sprintf("INSERT INTO %s (role_id, permission_id) VALUES (?, ?)", ds.tableName)
	_, err = ds.db.Exec(insertQuery, req.RoleID, req.PermissionID)
	if err != nil {
		return err
	}

	return nil

}

func (ds *RolePermissionDBDS) selectRolePermission(ctx context.Context, ID int64) (dataModel.RolePermission, error) {
	var rolePer dataModel.RolePermission
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE ID = ?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&rolePer.RoleID, &rolePer.PermissionID)
	if err != nil {
		return dataModel.RolePermission{}, err
	}
	return rolePer, nil
}
