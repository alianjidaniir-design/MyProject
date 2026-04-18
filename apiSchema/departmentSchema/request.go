package departmentSchema

type CreateDepartmentReq struct {
	College          string `json:"college"`
	EducationalGroup string `json:"educational_group"`
}

type ListReq struct {
	Page int `json:"page"`
	Size int `json:"size"`
}

type DeleteDepartmentReq struct {
	ID int64 `json:"id"`
}

type UpdateDepartmentReq struct {
	ID int64 `json:"id"`
}
