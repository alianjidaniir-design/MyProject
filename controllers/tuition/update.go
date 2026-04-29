package tuition

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/tuitionSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Update(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "16")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[tuitionSchema.UpdateTuition]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TuitionErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.TuitionRepo.Update(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TuitionErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
