package student

import (
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func MyInfo(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "13")
	defer mainController.FinishSpan(ctx)
	errStr, code, err := mainController.ParseBody(ctx, nil)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "06", errStr, code, err)
	}
	res, errStr, code, err := repositories.StudentRepo.MyInformation(spanCtx, ctx)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "08", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
