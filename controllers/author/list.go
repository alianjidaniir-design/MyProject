package author

import (
	"MyProject/apiSchema/authorSchema"
	"MyProject/apiSchema/commonSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func List(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "15")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[authorSchema.Pagination]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.AuthorErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.AuthorRepo.List(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.AuthorErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)

}
