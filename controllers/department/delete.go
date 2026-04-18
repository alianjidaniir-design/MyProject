package course

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/departmentSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Delete(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "12")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[departmentSchema.DeleteDepartmentReq]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.DepartmentErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.DepartmentRepo.Delete(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.DepartmentErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
