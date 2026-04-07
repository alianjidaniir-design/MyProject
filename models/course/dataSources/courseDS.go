package dataSources

import (
	"MyProject/apiSchema/courseSchema"
	courseDataModle "MyProject/models/course/dataModels"
	"context"
)

type CourseDB interface {
	CreateCourse(ctx context.Context, req courseSchema.RequestCourse) (courseDataModle.Course, error)
}
