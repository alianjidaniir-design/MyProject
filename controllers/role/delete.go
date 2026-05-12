package role

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/roleSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Delete(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "13")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[roleSchema.GetRoleRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.RoleErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.RoleRepo.Delete(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.RoleErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
