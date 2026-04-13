package teacherSchema

import "MyProject/models/teachers/dataModels"

type TeacherSchema struct {
	Teacher dataModels.Teacher `json:"teacher"`
}
