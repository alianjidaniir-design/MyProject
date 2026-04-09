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
	Get(ctx context.Context, req commonSchema.BaseRequest[courseSchema.GetCoursesRequest]) (res courseSchema.GetCoursesResponse, errStr string, code int, err error)
	Update(ctx context.Context, req commonSchema.BaseRequest[courseSchema.UpdateCourseRequest]) (res courseSchema.UpdateCourseResponse, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[courseSchema.HardDeleteCourseRequest]) (res courseSchema.HardDeleteCourseResponse, errStr string, code int, err error)
	SoftDelete(ctx context.Context, req commonSchema.BaseRequest[courseSchema.SoftDeleteCourseRequest]) (res courseSchema.SoftDeleteCourseResponse, errStr string, code int, err error)
}

var CourseRepo CourseRepository = course.GetRepoIns()
