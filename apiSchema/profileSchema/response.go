package profileSchema

import "MyProject/models/profile/dataModels"

type InformationResponse struct {
	Details dataModels.Profile `json:"details"`
	Message string             `json:"message"`
}

type ListAllScoresResp struct {
	Students []dataModels.ScoresStudents
	Total    int
}

type StudentsSummeryResponse struct {
	Summery []dataModels.StudentsSummary
	Total   int
}

type DetailProfileStudent struct {
	Detail []dataModels.ScoresAnnouncement `json:"detail"`
}

type DeleteProfileScoresResp struct {
	Message string `json:"message"`
}
