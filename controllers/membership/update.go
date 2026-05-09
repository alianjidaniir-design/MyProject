package membership

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/membershipSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Update(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "14")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[membershipSchema.UpdateMembership]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.MemberShipErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.MemberShipRepo.Update(spanCtx, req)
	if err != nil {
		return err
	}
	return mainController.Response(ctx, res)

}
