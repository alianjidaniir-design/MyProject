package teacher

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/teacherSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func SoftDelete(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "20")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[teacherSchema.SelectTeacherSchema]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TeacherErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.TeacherRepo.SoftDelete(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.TeacherErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
