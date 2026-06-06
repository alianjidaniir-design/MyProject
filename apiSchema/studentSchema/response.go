package studentSchema

import (
	"MyProject/models/student/dataModel"
)

type DetailStudent struct {
	Student dataModel.Student `json:"student"`
}

type InfoStudent struct {
	Info dataModel.InfoStudent `json:"info"`
}
type ListStudents struct {
	Students []dataModel.Student `json:"students"`
	Total    int64
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
	Massage string `json:"massage"`
}

type RefreshTokenResponse struct {
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}
