package authorSchema

import "MyProject/models/author/dataModel"

type DetailAuthor struct {
	Detail  dataModel.Author `json:"detail"`
	Massage string           `json:"massage"`
}
