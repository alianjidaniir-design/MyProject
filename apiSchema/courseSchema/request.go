package courseSchema

type RequestCourse struct {
	CourseNumber string `json:"course_number"`
	Title        string `json:"title"`
	Unit         int    `json:"unit"`
	DepartmentID int64  `json:"department_id"`
	Description  string `json:"description"`
}

type CoursesListRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type GetCoursesRequest struct {
	ID int64 `json:"ID"`
}

type UpdateCourseRequest struct {
	ID int64 `json:"ID"`
}

type HardDeleteCourseRequest struct {
	ID int64 `json:"ID"`
}

type SoftDeleteCourseRequest struct {
	ID int64 `json:"ID"`
}

type DeActiveCourseRequest struct {
	ID         int64 `json:"ID"`
	Deactivate bool  `json:"deactivate"`
}
