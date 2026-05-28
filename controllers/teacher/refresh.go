package teacher

import (
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Refresh(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "11")
	defer mainController.FinishSpan(ctx)
	errStr, code, err := mainController.ParseBody(ctx, nil)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TeacherErrCode, "01", errStr, code, err)
	}
	errStr, code, err = repositories.TeacherRepo.RefreshToken(spanCtx, ctx)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TeacherErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, nil)

}
