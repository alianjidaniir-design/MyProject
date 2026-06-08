package offeringSchema

import "time"

type CreateOfferingRequest struct {
	GroupNumber    int       `json:"group_number" validate:"required"`
	CourseNumber   int64     `json:"course_number" validate:"required"`
	TeacherId      int64     `json:"teacher_id" validate:"required"`
	Capacity       int       `json:"capacity" validate:"required"`
	IsActive       bool      `json:"is_active" validate:"omitempty"`
	TermId         int       `json:"term_id" validate:"required"`
	Week           string    `json:"week" validate:"required"`
	Day            string    `json:"day" validate:"required"`
	ClassStartTime string    `json:"class_start_time" validate:"required"`
	ClassEndTime   string    `json:"class_end_time" validate:"required"`
	ExamStartTime  time.Time `json:"exam_start_time" validate:"required"`
	ExamEndTime    time.Time `json:"exam_end_time" validate:"required"`
}

type ListOfferingsRequest struct {
	PageNumber int `json:"page_number"`
	PageSize   int `json:"page_size"`
}

type GetRowOfferingRequest struct {
	Row int64 `json:"row"`
}
