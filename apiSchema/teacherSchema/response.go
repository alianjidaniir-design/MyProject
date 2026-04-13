package teacherSchema

import "MyProject/models/teachers/dataModels"

type TeacherSchema struct {
	Teacher dataModels.Teacher `json:"teacher"`
}

type ListSchema struct {
	Teachers []dataModels.Teacher `json:"teachers"`
	Total    int64                `json:"total"`
}

type DetailTeacherSchema struct {
	Teacher dataModels.Teacher `json:"teacher"`
}

type HardDeleteTeacherSchema struct {
	Massage string `json:"massage"`
}

type SoftDeleteTeacherSchema struct {
	Teacher dataModels.Teacher `json:"teacher"`
	Massage string             `json:"massage"`
}
