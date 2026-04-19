package offeringSchema

import "MyProject/models/offering/dataModels"

type CreateOfferingResponse struct {
	Specification dataModels.Offering
}

type ListOfferingResponse struct {
	Offerings  []dataModels.Offering
	TotalCount int
}
