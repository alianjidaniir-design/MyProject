package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"MyProject/models/user"
	"context"
)

type UserRepository interface {
	// Create متد create
	Create(ctx context.Context, req commonSchema.BaseRequest[userSchema.LoginRequest]) (res userSchema.ResponseUser, errStr string, code int, err error)
	// List method list

	List(ctx context.Context, req commonSchema.BaseRequest[userSchema.ListRequest]) (res userSchema.ListUser, errStr string, code int, err error)

	// Get method
	Get(ctx context.Context, req commonSchema.BaseRequest[userSchema.GetRequest]) (res userSchema.GetResponse, errStr string, code int, err error)

	// Update method

	Update(ctx context.Context, req commonSchema.BaseRequest[userSchema.UpdateUserRequest]) (res userSchema.UpdateResponse, errStr string, code int, err error)
}

var UserRepo UserRepository = user.GetRepoIns()
