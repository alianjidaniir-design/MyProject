package dataSources

import (
	"MyProject/apiSchema/courseSchema"
	courseDataModle "MyProject/models/course/dataModels"
	"context"
)

type CourseDB interface {
	CreateCourse(ctx context.Context, req courseSchema.RequestCourse) (courseDataModle.Course, error)
	ListCourse(ctx context.Context, req courseSchema.CoursesListRequest) ([]courseDataModle.Course, int64, error)
	GetCourse(ctx context.Context, req courseSchema.GetCoursesRequest) (courseDataModle.Course, error)
	UpdateCourse(ctx context.Context, req courseSchema.UpdateCourseRequest) (courseDataModle.Course, error)
	DeleteCourse(ctx context.Context, req courseSchema.HardDeleteCourseRequest) (courseDataModle.Course, error)
	SoftDelete(ctx context.Context, req courseSchema.SoftDeleteCourseRequest) (courseDataModle.Course, error)
	ListDepartmentsCourse(ctx context.Context, req courseSchema.DepartmentListRequest) ([]courseDataModle.Course, int64, error)
}
