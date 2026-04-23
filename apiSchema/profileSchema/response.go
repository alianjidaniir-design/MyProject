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
