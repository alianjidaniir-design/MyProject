package user

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/userSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Get(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "12")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[userSchema.GetRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "40", errStr, code, err)
	}
	res, errStr, code, err := repositories.UserRepo.Get(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "41", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
