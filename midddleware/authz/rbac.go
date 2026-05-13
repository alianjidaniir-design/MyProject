package authz

import (
	dataModel2 "MyProject/models/permission/dataModel"
	"MyProject/models/role/dataModel"
	"context"
	"database/sql"
	"fmt"
)

func GetRoleByID(ds *sql.DB, RoleName string) (*dataModel.Role, error) {
	var role dataModel.Role
	var ctx context.Context
	selectQuery := "SELECT * FROM roles WHERE mame = ? "
	err := ds.QueryRowContext(ctx, selectQuery, RoleName).Scan(&role.ID, &RoleName)
	if err != nil {
		return nil, err
	} else if err == sql.ErrNoRows {
		return nil, fmt.Errorf("role '%s' not found", RoleName)
	}
	return &role, err
}

func listPermissions(db *sql.DB, roleID int64) ([]dataModel2.Permission, error) {
	var perm []dataModel2.Permission
	joinQuery := `
SELECT p.id , p.name FROM permission p
JOIN role_permission rp ON p.id = rp.permission_id
WHERE rp.role_id = ?`
	rows, err := db.Query(joinQuery, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var permission dataModel2.Permission
		err = rows.Scan(&permission.ID, &permission.Name)
		if err != nil {
			return nil, err
		}
		perm = append(perm, permission)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return perm, nil

}

func HasPermissionByTID(role *dataModel.Role, p dataModel2.Permission) (bool, error) {
	if role == nil {
		return false, nil
	}
	var db *sql.DB
	perm, err := listPermissions(db, role.ID)
	if err != nil {
		return false, err
	}
	for _, per := range perm {
		if per.Name == p.Name {
			return true, nil
		}
	}
	return false, err
}
