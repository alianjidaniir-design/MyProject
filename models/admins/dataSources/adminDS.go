package dataSources

import (
	"MyProject/apiSchema/adminSchema"
	adminDataModels "MyProject/models/admins/dataModels"
	"context"
)

type AdminDS interface {
	CreateAdmin(ctx context.Context, req adminSchema.InformationSchema) (res adminDataModels.Admins, err error)
	Login(ctx context.Context, req adminSchema.LoginAdminRequest) (access string, refresh string, massage string, err error)
	Refresh(ctx context.Context, req string) (access string, refresh string, err error)
	Logout(ctx context.Context, req adminSchema.LogoutSchema, ref string) error
}
