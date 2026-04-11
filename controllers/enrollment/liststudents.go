package enrollment

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/enrollmentSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func ListStudents(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "20")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[enrollmentSchema.ListCourseStudentsRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.CourseErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.EnrollmentRepos.ListCourseStudent(spanCtx, req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.CourseErrCode, "02", errStr, code, err)
	}
	return mainController.Response(ctx, res)
}
