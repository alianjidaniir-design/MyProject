package roleSchema

import "MyProject/models/role/dataModel"

type DetailRole struct {
	Role    dataModel.Role `json:"role"`
	Massage string         `json:"massage"`
}

type ListRole struct {
	DetailRole []dataModel.Role `json:"detailRole"`
	Total      int              `json:"total"`
}
