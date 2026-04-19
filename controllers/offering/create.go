package offering

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/offeringSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "11")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[offeringSchema.CreateOfferingRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.OfferingErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.OfferingRepo.Create(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.OfferingErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)

}
