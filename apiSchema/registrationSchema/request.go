package registrationSchema

type RegisterStudentRequest struct {
	StudentID int64             `json:"student_id" validate:"omitempty"`
	Selection []CourseSelection `json:"selection" validate:"required,min=1,dive"`
}

type CourseSelection struct {
	IsReserve    bool  `json:"is_reserve" validate:"omitempty"`
	CourseNumber int64 `json:"course_number" validate:"required"`
	OfferingID   int64 `json:"offering_row" validate:"required"`
}

type Pages struct {
	Term int `json:"term" validate:"required"`
	Year int `json:"year" validate:"required"`
	Page int `json:"page" validate:"omitempty"`
}

type GetRegisteredStudentsRequest struct {
	ID int64 `json:"id"`
}

type SelectPageRegisteredStudentsRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type ListStudentsRequest struct {
	StudentID int64 `json:"student_id"`
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
}

type ListOfferingRequest struct {
	OfferingRow int64 `json:"offering_row"`
	Page        int   `json:"page"`
	PageSize    int   `json:"page_size"`
}
