package translator

import (
	"MyProject/apiSchema/commonSchema"
	translateSchema "MyProject/apiSchema/translatorSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "12")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[translateSchema.CreateTranslator]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TransactionErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.TranslatorRepo.Create(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TransactionErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)

}
