package userSchema

import "MyProject/pkg/pagination"

type LoginRequest struct {
	Code   string `msgpack:"code" validate:"required , max = 10 "`
	Name   string `msgpack:"name" validate:"required"`
	Family string `msgpack:"family" validate:"required"`
}

type ListRequest struct {
	Page    pagination.Page    `json:"page" msgpack:"page"`
	PerPage pagination.PerPage `json:"perPage" msgpack:"perPage"`
}
