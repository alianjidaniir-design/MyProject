package user

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"MyProject/statics/constants/status"
	"context"
)

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.LoginRequest]) (res userSchema.ResponseUser, errStr string, code int, err error) {
	createdUser, err := repo.db().CreateUser(ctx, req.Body)
	if err != nil {
		return userSchema.ResponseUser{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.ResponseUser{User: createdUser}, "", status.StatusOK, nil
}
