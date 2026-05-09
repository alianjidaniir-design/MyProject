package membership

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/membershipSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Delete(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "13")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[membershipSchema.GetIDMembership]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.MemberShipErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.MemberShipRepo.Delete(spanCtx, req)
	if err != nil {
		return err
	}
	return mainController.Response(ctx, res)

}
