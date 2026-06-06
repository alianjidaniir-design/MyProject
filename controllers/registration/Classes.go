package registration

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/registrationSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Classes(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "18")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[registrationSchema.Pages]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.RegistrationErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.RegistrationRepo.ListClasses(spanCtx, req, ctx)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.RegistrationErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
