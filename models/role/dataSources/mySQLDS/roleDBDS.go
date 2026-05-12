package mySQLDS

import (
	"MyProject/apiSchema/roleSchema"
	"MyProject/models/role/dataModel"
	"MyProject/models/role/dataSources"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type RoleDBDS struct {
	tableName string
	db        *sql.DB
}

func NewRoleDBDS(tableName string, db *sql.DB) (dataSources.RoleDS, error) {
	ff := &RoleDBDS{
		tableName: tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *RoleDBDS) CreateRole(ctx context.Context, req roleSchema.CreateRoleRequest) (res dataModel.Role, err error) {
	if req.Name == "" || len(req.Name) > 32 {
		return dataModel.Role{}, errors.New("name length must be between 32 characters and name is not empty")
	}
	var tot int
	countName := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name=?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countName).Scan(&tot)
	if err != nil {
		return dataModel.Role{}, err
	} else if tot > 1 {
		return dataModel.Role{}, errors.New("there is already a role with this name")
	}
	insertQuery := fmt.Sprintf("INSERT INTO `%s` (`name`) VALUES (?)", ds.tableName)
	lastID, err := ds.db.ExecContext(ctx, insertQuery, req.Name)
	if err != nil {
		return dataModel.Role{}, err
	}
	newID, err := lastID.LastInsertId()
	if err != nil {
		return dataModel.Role{}, err
	}
	return ds.selectRoles(ctx, newID)

}

func (ds *RoleDBDS) selectRoles(ctx context.Context, ID int64) (dataModel.Role, error) {
	var role dataModel.Role
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id=?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&role.ID, &role.Name)
	if err != nil {
		return role, err
	}
	return role, nil
}
