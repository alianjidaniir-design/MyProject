package authz

import (
	dataModel2 "MyProject/models/permission/dataModel"
	"MyProject/models/role/dataModel"
	"database/sql"
	"errors"
	"fmt"
)

func GetRoleByID(RoleName string) (*dataModel.Role, error) {
	var role dataModel.Role
ro , err:=
}

func listPermissions(roleID int64) ([]dataModel2.Permission, error) {
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

func HasPermissionByTID(role *dataModel.Role, p string) (bool, error) {
	if role == nil {
		return false, nil
	}
	perm, err := listPermissions(role.ID)
	if err != nil {
		return false, err
	}
	for _, per := range perm {
		if per.Name == p {
			return true, nil
		}
	}
	return false, err
}
