package userSchema

import "MyProject/pkg/pagination"

type LoginRequest struct {
	Code   string `json:"code"  validate:"required , max = 10"`
	Name   string `json:"name" validate:"required"`
	Family string `json:"family" validate:"required"`
}

type ListRequest struct {
	Page    pagination.Page    `json:"page"`
	PerPage pagination.PerPage `json:"perPage"`
}

type GetRequest struct {
	ID int64 `json:"id" validate:"required"`
}
