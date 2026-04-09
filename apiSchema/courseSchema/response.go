package courseSchema

import (
	courseDataModle "MyProject/models/course/dataModels"
)

type ResponseCourse struct {
	Course courseDataModle.Course `json:"course"`
}

type CourseListResponse struct {
	Courses []courseDataModle.Course `json:"courses"`
	Total   int64                    `json:"total"`
}

type GetCoursesResponse struct {
	Courses courseDataModle.Course `json:"courses"`
}

type UpdateCourseResponse struct {
	Course courseDataModle.Course `json:"course"`
}

type HardDeleteCourseResponse struct {
	Course courseDataModle.Course `json:"course"`
}

type SoftDeleteCourseResponse struct {
	Course courseDataModle.Course `json:"course"`
}

type DeactivateCourseResponse struct {
	Massage string `json:"massage"`
}
