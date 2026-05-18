package repositories

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/subjectSchema"
	"MyProject/models/subject"
	"context"
)

type SubjectRepository interface {
	Create(ctx context.Context, req commonSchema.BaseRequest[subjectSchema.CreateSubject]) (res subjectSchema.DetailSubject, errStr string, code int, err error)
	Get(ctx context.Context, req commonSchema.BaseRequest[subjectSchema.GetSubject]) (res subjectSchema.DetailSubject, errStr string, code int, err error)
	Delete(ctx context.Context, req commonSchema.BaseRequest[subjectSchema.GetSubject]) (res subjectSchema.DetailSubject, errStr string, code int, err error)
	List(ctx context.Context , req commonSchema.BaseRequest[subjectSchema.Pagination]) (res subjectSchema.ListSubjects, errStr string, code int, err error)
}

var SubjectRepo SubjectRepository = subject.GetRepo()
