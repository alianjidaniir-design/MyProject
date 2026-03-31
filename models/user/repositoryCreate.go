package user

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"MyProject/statics/constants/status"
	"context"
	"errors"
)

func (repo *Repository) Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.LoginRequest]) (res userSchema.ResponseUser, errStr string, code int, err error) {
	if repo.initErr != nil {
		return userSchema.ResponseUser{}, "13", status.StatusUnauthorized, repo.initErr
	}
	if repo.db() == nil {
		return userSchema.ResponseUser{}, "14", status.UnAvailableServiceError, errors.New("student datasourse not configured")
	}

	createdUser, err := repo.db().CreateStudent(ctx, req.Body)
	if err != nil {
		return userSchema.ResponseUser{}, "04", status.UnAvailableServiceError, err
	}
	return userSchema.ResponseUser{User: createdUser}, "", status.StatusOK, nil
}
