package subject

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/subjectSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Get(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "1۳")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[subjectSchema.GetSubject]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.SubjectErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.SubjectRepo.Get(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.SubjectErrCode, "03", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
