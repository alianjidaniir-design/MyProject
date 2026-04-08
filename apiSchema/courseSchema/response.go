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
