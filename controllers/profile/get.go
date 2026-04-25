package profile

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/profileSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Get(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "15")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[profileSchema.GetScoresReq]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.ProfileErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.ProfileRepo.Get(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.ProfileErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)

}
