package user

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"MyProject/statics/constants/status"
	"context"
	"errors"
)

func (repo *Repository) List(ctx context.Context, req commonSchema.BaseRequest[userSchema.ListRequest]) (res userSchema.ListUser, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.ListUser{}, "10", status.UnAvailableServiceError, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.ListUser{}, "11", status.StatusInternalServerError, errors.New("bad")
	}
	ListUsers, err := repo.db().ReadStudent(ctx, req.Body)
	if err != nil {
		return userSchema.ListUser{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.ListUser{Users: []ListUsers}, "", status.StatusOK, nil
}
