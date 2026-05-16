package authorSchema

import "MyProject/models/author/dataModel"

type DetailAuthor struct {
	Detail  dataModel.Author `json:"detail"`
	Massage string           `json:"massage"`
}

type ListAuthor struct {
	Massage string             `json:"massage"`
	List    []dataModel.Author `json:"list"`
	Total   int                `json:"total"`
}
