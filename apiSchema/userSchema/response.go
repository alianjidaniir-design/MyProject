package userSchema

import (
	UserdataModel "MyProject/models/user/dataModel"
)

type ResponseUser struct {
	User UserdataModel.User `json:"user"`
}

type ListUser struct {
	Users []UserdataModel.User
	Total int64
}

type GetResponse struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Family string `json:"family"`
}
