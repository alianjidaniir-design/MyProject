package categorySchema

import "MyProject/models/category/dataModel"

type InformationCategoryResponse struct {
	Detail  dataModel.Category `json:"detail"`
	Massage string             `json:"massage"`
}

type ListCategoryResponse struct {
	List  []dataModel.Category `json:"list"`
	Total int                  `json:"total"`
}
