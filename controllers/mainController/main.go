package mainController

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/statics/constants/status"
	"context"
	"fmt"
	"reflect"

	"github.com/gofiber/fiber/v2"
)

type Errorresponse struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

type responseUser struct {
	Data any `json:"data"`
}

func InitAPI(ctx *fiber.Ctx, sectionErrCode string) context.Context {
	_ = ctx
	_ = sectionErrCode
	return context.Background()
}

func FinishSpan(ctx *fiber.Ctx) {
	_ = ctx
}

func ParseBody(ctx *fiber.Ctx, req any) (string, int, error) {
	if err := ctx.BodyParser(req); err != nil {
		return "01", status.StatusBadRequest, err
	}
	fillHeaders(ctx, req)

	return "", status.StatusOK, nil
}

func ParseQuery(ctx *fiber.Ctx, req any) (string, int, error) {
	if err := ctx.QueryParser(req); err != nil {
		return "02", status.StatusBadRequest, err
	}
	headers := map[string]string{}
	for k, v := range ctx.GetReqHeaders() {
		headers[k] = v[0]
	}
	validator, ok := req.(interface {
		Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error)
	})
	if !ok {
		return "", status.StatusBadRequest, nil
	}
	return validator.Validate(commonSchema.ValidateExtraData{Headers: headers})
}

func Error(ctx *fiber.Ctx, baseErrCode string, section string, errStr string, code int, err error) error {
	return ctx.Status(code).JSON(Errorresponse{
		Code:    fmt.Sprintf("%s%s%s", baseErrCode, section, errStr),
		Message: err.Error(),
	})
}

func Response(ctx *fiber.Ctx, res any) error {
	return ctx.Status(status.StatusOK).JSON(responseUser{Data: res})
}

func fillHeaders(ctx *fiber.Ctx, req any) {
	refVal := reflect.ValueOf(req)
	if refVal.Kind() != reflect.Ptr || refVal.Elem().Kind() != reflect.Struct {
		return
	}
	headersfield := refVal.Elem().FieldByName("Headers")
	if !headersfield.IsValid() || !headersfield.CanSet() || headersfield.Kind() != reflect.Map {
		return
	}
	header := map[string]string{}
	for key, val := range ctx.GetReqHeaders() {
		header[key] = val[0]
	}
	headersfield.Set(reflect.ValueOf(header))

}
