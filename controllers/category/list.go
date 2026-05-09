package category

import (
	"MyProject/apiSchema/categorySchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func List(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "15")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[categorySchema.PaginationList]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.CategoryErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.CategoryRepo.List(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.CategoryErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
