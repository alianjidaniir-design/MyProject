package dataSourses

import (
	"MyProject/apiSchema/userSchema"
	userDataModel "MyProject/models/user/dataModel"
	"context"
)

type UserDBDS interface {
	CreateStudent(ctx context.Context, req userSchema.LoginRequest) (userDataModel.User, error)
}
