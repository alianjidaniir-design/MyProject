package category

import (
	"MyProject/apiSchema/categorySchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Get(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "14")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[categorySchema.GetRowCategoryRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.CategoryErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.CategoryRepo.Get(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.CategoryErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
