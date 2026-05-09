package programSchema

import "MyProject/models/program/dataModel"

type DetailProgramResponse struct {
	Detail  dataModel.Program `json:"detail"`
	Massage string            `json:"massage"`
}

type ListProgramsResponse struct {
	Programs []dataModel.Program `json:"programs"`
	Total    int64               `json:"total"`
}
