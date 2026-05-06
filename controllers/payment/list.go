package payment

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/paymentSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func List(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "14")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[paymentSchema.ListPayment]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.PaymentErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.PaymentRepo.List(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.PaymentErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)

}
