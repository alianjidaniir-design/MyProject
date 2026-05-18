package dataSources

import (
	"MyProject/apiSchema/subjectSchema"
	"MyProject/models/subject/dataModel"
	"context"
)

type SubjectDS interface {
	CreateSubject(ctx context.Context, req subjectSchema.CreateSubject) (res dataModel.Subject, err error)
	GetSubject(ctx context.Context, req subjectSchema.GetSubject) (res dataModel.Subject, err error)
	DeleteSubject(ctx context.Context, req subjectSchema.GetSubject) (res dataModel.Subject, err error)
	ListSubjects(ctx context.Context, req subjectSchema.Pagination) (res []dataModel.Subject,total int, err error)
}
