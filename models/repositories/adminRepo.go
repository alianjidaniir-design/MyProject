package repositories

import (
	"MyProject/apiSchema/adminSchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/models/admins"
	"context"

	"github.com/gofiber/fiber/v2"
)

type AdminRepository interface {
	// Create method
	Create(ctx context.Context, req commonSchema.BaseRequest[adminSchema.InformationSchema]) (res adminSchema.DetailAdminSchema, errStr string, code int, err error)
	Login(ctx context.Context, req commonSchema.BaseRequest[adminSchema.LoginAdminRequest], c *fiber.Ctx) (res adminSchema.EntryAdminSchema, errStr string, code int, err error)
	Refresh(ctx context.Context, c *fiber.Ctx) (errStr string, code int, err error)
	Logout(ctx context.Context, req commonSchema.BaseRequest[adminSchema.LogoutSchema], c *fiber.Ctx) (res adminSchema.EntryAdminSchema, errStr string, code int, err error)
}

var AdminRepo AdminRepository = admins.GetRepo()
