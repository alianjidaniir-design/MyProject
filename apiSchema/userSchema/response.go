package userSchema

import UserdataModel "MyProject/models/user/dataModel"

type ResponseUser struct {
	User UserdataModel.User
}

type ListUser struct {
	Users []UserdataModel.User
}
