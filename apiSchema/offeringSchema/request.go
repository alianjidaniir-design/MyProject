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
	Term             int    `json:"term" validate:"required"`
	Year             int    `json:"year" validate:"required"`
	College          string `json:"college" validate:"omitempty"`
	EducationalGroup string `json:"educational_group" validate:"omitempty"`
	Week             string `json:"week" validate:"omitempty"`
	Day              string `json:"day" validate:"omitempty"`
}

type GetRowOfferingRequest struct {
	Row int64 `json:"row"`
}

type EditOffering struct {
	Row            int64  `json:"row" validate:"required"`
	Capacity       int    `json:"capacity" validate:"omitempty"`
	IsActive       *bool  `json:"is_active" validate:"omitempty"`
	ExamStartTime  string `json:"exam_start_time" validate:"omitempty"`
	ExamFinishTime string `json:"exam_finish_time" validate:"omitempty"`
}

type Pages struct {
	Term int `json:"term" validate:"required"`
	Year int `json:"year" validate:"required"`
	Page int `json:"page" validate:"required"`
}
