package user

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func List(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "12")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[userSchema.ListRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "06", errStr, code, err)
	}
	res, errStr, code, err := repositories.UserRepo.List(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "08", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
