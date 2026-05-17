package publisher

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/publisherSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "12")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[publisherSchema.CreatePublisher]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.PublisherErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.PublisherRepo.Create(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.PublisherErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
