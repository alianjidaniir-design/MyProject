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
	User UserdataModel.User `json:"user"`
}

type UpdateResponse struct {
	User UserdataModel.User `json:"user"`
}
