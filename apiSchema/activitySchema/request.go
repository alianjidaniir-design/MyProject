package activitySchema

type CreateActivity struct {
	ActivityType string `json:"activityType" validate:"required,oneOf= exercise, practice, assignment"`
	Title        string `json:"title" validate:"required,max=130"`
	Description  string `json:"description" validate:"required"`
	TeacherID    string `json:"teacherID" validate:"required"`
	CourseID     string `json:"courseID" validate:"required"`
	GroupNumber  int    `json:"groupNumber" validate:"required,numeric"`
	IsActive     bool   `json:"isActive" validate:"omitempty,default true"`
}
