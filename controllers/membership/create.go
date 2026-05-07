package membership

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/membershipSchema"
	"MyProject/controllers/mainController"
	"MyProject/models/repositories"
	"MyProject/statics/constants/controllerbaseErrCode"

	"github.com/gofiber/fiber/v2"
)

func Create(ctx *fiber.Ctx) error {
	spanCtx := mainController.InitAPI(ctx, "12")
	defer mainController.FinishSpan(ctx)
	req := commonSchema.BaseRequest[membershipSchema.CreateMembershipRequest]{}
	errStr, code, err := mainController.ParseBody(ctx, &req)
	if err != nil {
		return mainController.Error(ctx, controllerbaseErrCode.MemberShipErrCode, "01", errStr, code, err)
	}
	res, errStr, code, err := repositories.MemberShipRepo.Create(spanCtx, req)
	if err != nil {
		return err
	}
	return mainController.Response(ctx, res)

}
