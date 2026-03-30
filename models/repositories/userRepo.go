package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.LoginRequest]) (res userSchema.LoginRequest, errStr string, code int, err error)
}

var UserRepo UserRepository
