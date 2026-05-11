package userSchema

import (
	"MyProject/models/user/dataModel"
)

type ResponseUser struct {
	User dataModel.User `json:"user"`
}

type ListUser struct {
	Users []dataModel.User `json:"users"`
	Total int64
}

type GetResponse struct {
	User dataModel.User `json:"user"`
}

type UpdateResponse struct {
	User dataModel.User `json:"user"`
}

type DeleteResponse struct {
	User dataModel.User `json:"user"`
}

type SoftDeleteResponse struct {
	User dataModel.User `json:"user"`
}
