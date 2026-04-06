package courseSchema

type RequestCourse struct {
	CourseCode string `json:"course_code"`
	Title      string `json:"title"`
	Capacity   int    `json:"capacity"`
	IsActive   bool   `json:"isActive"`
}
