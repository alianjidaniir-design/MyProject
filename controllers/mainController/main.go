package mainController

import (
	"MyProject/apiSchema/commonSchema"

	"github.com/gofiber/fiber/v2"
)

type request interface {
	Validate(validateExtraData commonSchema.ValidateExtraData) (string, int, error)
}

func ParseBody[T any](ctx *fiber.Ctx, request *commonSchema.BaseRequest[T]) (string, int, error) {
	var (
		err        error
		codeStr    string
		statusCode int
		body       = ctx.Body()
	)
	request.FillHeaders(ctx)
	if len(body) != 0 {

	}
}
