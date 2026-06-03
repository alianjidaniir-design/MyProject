package courseSchema

type RequestCourse struct {
	CourseNumber int64  `json:"course_number" validate:"required"`
	Title        string `json:"title" validate:"required,max=64"`
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
	ID int64 `json:"id"`
}

type UpdateCourseRequest struct {
	ID int64 `json:"id"`
}

type HardDeleteCourseRequest struct {
	ID int64 `json:"id"`
}

type SoftDeleteCourseRequest struct {
	ID int64 `json:"id"`
}

type DeActiveCourseRequest struct {
	ID         int64 `json:"id"`
	Deactivate bool  `json:"deactivate"`
}
