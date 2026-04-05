package userSchema

import "MyProject/pkg/pagination"

type LoginRequest struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Family string `json:"family"`
}

type ListRequest struct {
	Page    pagination.Page    `json:"page"`
	PerPage pagination.PerPage `json:"perPage"`
}

type GetRequest struct {
	ID int64 `json:"id"`
}
