package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/studentSchema"
	"MyProject/models/student"
	"context"
)

type StudentRepository interface {
	// Create متد create
	Create(ctx context.Context, req commonSchema.BaseRequest[studentSchema.SignUpStudent]) (res studentSchema.UserResponse, errStr string, code int, err error)
	// List method list

	List(ctx context.Context, req commonSchema.BaseRequest[studentSchema.ListRequest]) (res studentSchema.ListUser, errStr string, code int, err error)

	// Get method
	Get(ctx context.Context, req commonSchema.BaseRequest[studentSchema.GetRequest]) (res studentSchema.GetResponse, errStr string, code int, err error)

	// Update method

	Update(ctx context.Context, req commonSchema.BaseRequest[studentSchema.UpdateUserRequest]) (res studentSchema.UpdateResponse, errStr string, code int, err error)

	// Delete method

	Delete(ctx context.Context, req commonSchema.BaseRequest[studentSchema.DeleteRequest]) (res studentSchema.DeleteResponse, errStr string, code int, err error)

	SoftDelete(ctx context.Context, req commonSchema.BaseRequest[studentSchema.SoftDeleteRequest]) (res studentSchema.SoftDeleteResponse, errStr string, code int, err error)

	Entry(ctx context.Context, req commonSchema.BaseRequest[studentSchema.LoginStudent]) (res studentSchema.StudentEntry, errStr string, code int, err error)

	RefreshToken(ctx context.Context, req commonSchema.BaseRequest[studentSchema.RefreshTokenRequest]) (res studentSchema.RefreshTokenResponse, errStr string, code int, err error)
}

var StudentRepo StudentRepository = student.GetRepoIns()
