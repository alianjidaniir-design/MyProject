package bookSchema

import "MyProject/models/book/dataModel"

type InformationBook struct {
	Information dataModel.Book `json:"information"`
	Massage     string         `json:"massage"`
}

type ListBooks struct {
	Books []dataModel.Book `json:"books"`
	Total int              `json:"total"`
}
