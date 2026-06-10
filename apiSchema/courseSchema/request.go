package courseSchema

type RequestCourse struct {
	CourseNumber int64  `json:"course_number" validate:"required"`
	Title        string `json:"title" validate:"required,max=64"`
	CourseType   string `json:"course_type" validate:"required"`
	Unit         int    `json:"unit" validate:"required,gte=1,lte=9"`
	DepartmentID int64  `json:"department_id" validate:"required"`
	Prerequisite string `json:"prerequisite" validate:"omitempty,max=128"`
	Necessary    string `json:"necessary" validate:"omitempty,max=128"`
}

type CoursesListRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type DepartmentListRequest struct {
	DepartmentID int64 `json:"department_id"`
	Page         int   `json:"page"`
	PageSize     int   `json:"page_size"`
}

type GetCoursesRequest struct {
	CourseNumber int64 `json:"course_number"`
}

type UpdateCourseRequest struct {
	CourseNumber int64  `json:"course_number" validate:"omitempty"`
	NewCourseNum int64  `json:"new_course_num" validate:"omitempty"`
	Title        string `json:"title" validate:"omitempty,max=63"`
	CourseType   string `json:"course_type" validate:"omitempty"`
	Unit         int    `json:"unit" validate:"omitempty,gte=1,lte=9"`
	DepartmentID int64  `json:"department_id" validate:"omitempty"`
	Prerequisite string `json:"prerequisite" validate:"omitempty,max=128"`
	Necessary    string `json:"necessary" validate:"omitempty,max=128"`
}

type HardDeleteCourseRequest struct {
	CourseNumber int64 `json:"course_number"`
}

type SoftDeleteCourseRequest struct {
	CourseNumber int64 `json:"course_number"`
}
