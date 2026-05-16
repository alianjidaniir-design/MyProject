package authz

import (
	"MyProject/models/role/dataModel"
	"context"
	"errors"
	"strings"
)

func (a *AuthzMiddleWare) ParseRole(s string) (*dataModel.Role, error) {
	ctx := context.Background()
	towerName := strings.ToLower(s)
	if a.RolsDS == nil {
		return nil, errors.New("role data source is not initialized")
	}
	getRole, err := a.RolsDS.GetRoleByName(ctx, towerName)
	if err != nil {
		return nil, err
	}
	return &getRole, nil

}
