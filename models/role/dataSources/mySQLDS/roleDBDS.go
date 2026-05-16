package mySQLDS

import (
	"MyProject/apiSchema/roleSchema"
	"MyProject/models/role/dataModel"
	"MyProject/models/role/dataSources"
	"MyProject/pkg/pagination"
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
	countName := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE name = ?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countName, req.Name).Scan(&tot)
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

func (ds *RoleDBDS) GetRole(ctx context.Context, req roleSchema.GetRoleRequest) (res dataModel.Role, err error) {
	err = ds.check(ctx, req.ID)
	if err != nil {
		return dataModel.Role{}, err
	}
	return ds.selectRoles(ctx, req.ID)

}

func (ds *RoleDBDS) DeleteRole(ctx context.Context, req roleSchema.GetRoleRequest) (res dataModel.Role, err error) {
	err = ds.check(ctx, req.ID)
	if err != nil {
		return res, err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return res, err
	}
	return res, nil
}

func (ds *RoleDBDS) ListRoles(ctx context.Context, req roleSchema.Pagination) (res []dataModel.Role, total int, err error) {
	var roles []dataModel.Role
	page, pageSize, err := pagination.CheckPage(req.Page, req.Size)
	if err != nil {
		return nil, 0, err
	}
	limit := pageSize
	offset := (page - 1) * limit
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var role dataModel.Role
		err = rows.Scan(&role.ID, &role.Name)
		if err != nil {
			return nil, 0, err
		}
		roles = append(roles, role)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return roles, totalRows, nil

}

func (ds *RoleDBDS) check(ctx context.Context, ID int64) error {
	var check bool
	checking := `
SELECT EXISTS(SELECT 1 FROM ` + ds.tableName + ` WHERE ID = ?)`
	err := ds.db.QueryRowContext(ctx, checking, ID).Scan(&check)
	if err != nil {
		return err
	}
	if !check {
		return errors.New("role does not exist")
	}
	return nil
}

func (ds *RoleDBDS) GetRoleByName(ctx context.Context, name string) (dataModel.Role, error) {
	var role dataModel.Role
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE name = ?", ds.tableName)
	err := ds.db.QueryRowContext(ctx, query, name).Scan(&role.ID, &role.Name)
	if err != nil {
		return role, err
	}
	return role, nil
}
