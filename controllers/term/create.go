package term

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/termSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "13")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[termSchema.CreateTerm]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TermErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.TermRepo.Create(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TermErrCode, "03", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
