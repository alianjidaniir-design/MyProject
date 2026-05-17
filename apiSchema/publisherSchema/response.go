package publisherSchema

import "MyProject/models/publisher/dataModel"

type DetailPublisher struct {
	Massage string              `json:"massage"`
	Detail  dataModel.Publisher `json:"detail"`
}

type ListPublisherDetail struct {
	Massage string                `json:"massage"`
	List    []dataModel.Publisher `json:"list"`
	Total   int                   `json:"total"`
}
