package userSchema

import (
	UserdataModel "MyProject/models/user/dataModel"
	"MyProject/pkg/pagination"
)

type ResponseUser struct {
	User UserdataModel.User `json:"user"`
}

type ListUser struct {
	Users   []UserdataModel.User
	page    pagination.Page
	PerPage pagination.PerPage
}

type GetResponse struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Family string `json:"family"`
}
