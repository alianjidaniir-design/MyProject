package payment

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/paymentSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Get(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "12")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[paymentSchema.GetInformation]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.PaymentErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.PaymentRepo.Get(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.PaymentErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)

}
