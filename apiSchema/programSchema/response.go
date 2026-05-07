package programSchema

import "MyProject/models/program/dataModel"

type DetailProgramResponse struct {
	Detail  dataModel.Program `json:"detail"`
	Massage string            `json:"massage"`
}
