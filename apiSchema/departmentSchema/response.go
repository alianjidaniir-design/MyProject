package departmentSchema

import (
	"MyProject/models/department/dataModels"
)

type InformationDepartmentResp struct {
	Department dataModels.Department `json:"department"`
}

type ListDepartmentResp struct {
	Department []dataModels.Department `json:"department"`
	Total      int                     `json:"total"`
}

type UpdateDepartmentResp struct {
	Department dataModels.Department `json:"department"`
}
