package commonSchema

import (
	userModel "MyProject/models/user/dataModel"
	"MyProject/statics/constants/headers"
	"context"

	"github.com/gofiber/fiber/v2"
)

type Response struct {
	Status string      `json:"status" msgpack:"status"`
	Error  Error       `json:"error" msgpack:"error"`
	Data   interface{} `json:"data" msgpack:"data"`
}

type BaseRequest[T any] struct {
	SpanCtx context.Context
	User    userModel.User
	Req     T
	Headers
}

type Error struct {
	Message string `json:"message" msgpack:"message"`
	Code    int    `json:"code" msgpack:"code"`
}

type ValidateExtraData struct {
	Headers *Headers
}

type Headers struct {
	DUID string
	IP   string
}

func (header *Headers) FillHeaders(ctx *fiber.Ctx) {
	header.DUID = ctx.Get(headers.DUID)
}
