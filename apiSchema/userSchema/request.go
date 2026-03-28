package userSchema

type LoginRequest struct {
	Code   string `msgpack:"code" validate:"required , max = 10 "`
	Name   string `msgpack:"name" validate:"required"`
	Family string `msgpack:"family" validate:"required"`
}
