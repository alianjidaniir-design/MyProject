package student

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/studentSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Update(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "16")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[studentSchema.UpdateUserRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "06", errStr, code, err)
	}
	res, errStr, code, err := repositories.StudentRepo.Update(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.UserErrCode, "08", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
