package authz

import (
	dataSources2 "MyProject/models/permission/dataSources"
	"MyProject/models/role/dataModel"
	"MyProject/models/role/dataSources"
)

type AuthzMiddleWare struct {
	RolsDS       dataSources.RoleDS
	PermissionDS dataSources2.PermissionDS
}

func NewAuthMiddleware(roleDS dataSources.RoleDS, permDS dataSources2.PermissionDS) *AuthzMiddleWare {
	return &AuthzMiddleWare{RolsDS: roleDS, PermissionDS: permDS}
}

func (a *AuthzMiddleWare) HasPermissionByTID(role *dataModel.Role, p string) (bool, error) {
	if role == nil {
		return false, nil
	}
	perm, err := a.PermissionDS.ListPerms(role.ID)
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
