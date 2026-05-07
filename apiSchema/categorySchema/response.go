package categorySchema

import "MyProject/models/category/dataModel"

type InformationCategoryResponse struct {
	Detail  dataModel.Category `json:"detail"`
	Massage string             `json:"massage"`
}
