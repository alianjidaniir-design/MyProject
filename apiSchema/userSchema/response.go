package userSchema

import (
	UserdataModel "MyProject/models/user/dataModel"
	"MyProject/pkg/pagination"
)

type ResponseUser struct {
	User UserdataModel.User
}

type ListUser struct {
	Users   []UserdataModel.User
	page    pagination.Page
	PerPage pagination.PerPage
}
