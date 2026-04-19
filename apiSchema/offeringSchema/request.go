package offeringSchema

import "time"

type CreateOfferingRequest struct {
	GroupNumber    int       `json:"group_number"`
	CourseId       int64     `json:"course_id"`
	TeacherId      int64     `json:"teacher_id"`
	Capacity       int       `json:"capacity"`
	IsActive       bool      `json:"is_active"`
	TermId         int       `json:"term_id"`
	ClassStartTime time.Time `json:"class_start_time"`
	ClassEndTime   time.Time `json:"class_end_time"`
	ExamStartTime  time.Time `json:"exam_start_time"`
	ExamEndTime    time.Time `json:"exam_end_time"`
}

type ListOfferingsRequest struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

type GetRowOfferingRequest struct {
	Row int64 `json:"row"`
}
