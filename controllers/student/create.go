package student

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/studentSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "11")

	defer mainController.FinishSpan(ctx)

	req := commonSchema.BaseRequest[studentSchema.SignUpStudent]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "02", errStr, code, err)
	}
	res, errStr, code, err := repositories.StudentRepo.Create(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "03", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
