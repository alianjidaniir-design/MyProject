package studentSchema

import (
	"MyProject/models/student/dataModel"
)

type UserResponse struct {
	User dataModel.Student `json:"student"`
}

type ResponseUser struct {
	User dataModel.Student `json:"student"`
}

type ListUser struct {
	Users []dataModel.Student `json:"users"`
	Total int64
}

type GetResponse struct {
	User dataModel.Student `json:"student"`
}

type UpdateResponse struct {
	User dataModel.Student `json:"student"`
}

type DeleteResponse struct {
	User dataModel.Student `json:"student"`
}

type SoftDeleteResponse struct {
	User dataModel.Student `json:"student"`
}

type StudentEntry struct {
	Massage      string `json:"massage"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	RefreshToken dataModel.RefreshToken `json:"refresh_token"`
}
