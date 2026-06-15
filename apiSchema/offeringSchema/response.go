package offeringSchema

import "MyProject/models/offering/dataModels"

type CreateOfferingResponse struct {
	Specification dataModels.Offering
}

type ListOfferingResponse struct {
	Offerings  []dataModels.ListOfferings
	TotalCount int
}

type DetailOfferingResponse struct {
	Specification dataModels.Offering
}

type DeactivateOfferingResponse struct {
	Specification dataModels.Offering
	Massage       string `json:"massage"`
}

type ViewAfterEditCourse struct {
	Massage       string `json:"massage"`
	Specification dataModels.Offering
}

type ClassesTeacher struct {
	MyClasses []dataModels.TermClasses
	Page      int `json:"page"`
	Total     int
}
