package publisherSchema

import "MyProject/models/publisher/dataModel"

type DetailPublisher struct {
	Massage string              `json:"massage"`
	Detail  dataModel.Publisher `json:"detail"`
}
