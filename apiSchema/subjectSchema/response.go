package subjectSchema

import "MyProject/models/subject/dataModel"

type DetailSubject struct {
	Massage string            `json:"massage"`
	Detail  dataModel.Subject `json:"detail"`
}

type ListSubjects struct {
	Massage string `json:"massage"`
	List    []dataModel.Subject `json:"list"`
	Total   int                 `json:"total"`
}
