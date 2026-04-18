package termSchema

import "MyProject/models/term/dataModels"

type InformationTerm struct {
	Term dataModels.Term `json:"term"`
}

type ListTerms struct {
	Term  []dataModels.Term `json:"term"`
	Total int               `json:"total"`
}

type DeletedTerm struct {
	Message string `json:"message"`
}
