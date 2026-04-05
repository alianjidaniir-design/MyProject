package dataSourses

import (
	"MyProject/apiSchema/userSchema"
	userDataModel "MyProject/models/user/dataModel"
	"context"
)

type UserDB interface {
	CreateStudent(ctx context.Context, req userSchema.LoginRequest) (userDataModel.User, error)
	ReadStudent(ctx context.Context, req userSchema.ListRequest) ([]userDataModel.User, error)
}
