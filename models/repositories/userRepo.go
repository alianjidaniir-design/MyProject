package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"context"
)

type UserRepository interface {
	// Create متد create
	// List method list
	Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.LoginRequest]) (res userSchema.ResponseUser, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[userSchema.ListRequest]) (res userSchema.ListUser, errStr string, code int, err error)
}

var UserRepo UserRepository
