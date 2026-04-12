package courseSchema

type RequestCourse struct {
	CourseCode string `json:"course_code"`
	Title      string `json:"title"`
	Capacity   int    `json:"capacity"`
	IsActive   bool   `json:"isActive"`
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
