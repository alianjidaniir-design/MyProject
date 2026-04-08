package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/courseSchema"
	"MyProject/models/course"
	"context"
)

type CourseRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[courseSchema.RequestCourse]) (res courseSchema.ResponseCourse, errStr string, code int, err error)
	List(ctx context.Context, req commonSchema.BaseRequest[courseSchema.CoursesListRequest]) (res courseSchema.CourseListResponse, errStr string, code int, err error)
}

var CourseRepo CourseRepository = course.GetRepoIns()
