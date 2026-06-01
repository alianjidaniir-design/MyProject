package courseSchema

type RequestCourse struct {
	CourseNumber int64  `json:"course_number"`
	Title        string `json:"title"`
	Unit         int    `json:"unit"`
	DepartmentID int64  `json:"department_id"`
	Description  string `json:"description"`
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
	CourseNumber int64 `json:"course_number"`
}

type HardDeleteCourseRequest struct {
	CourseNumber int64 `json:"course_number"`
}

type SoftDeleteCourseRequest struct {
	CourseNumber int64 `json:"course_number"`
}

type DeActiveCourseRequest struct {
	CourseNumber int64 `json:"course_number"`
	Deactivate   bool  `json:"deactivate"`
}
