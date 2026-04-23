package registrationSchema

type RegisterStudentRequest struct {
	StudentID  int64 `json:"student_id"`
	OfferingID int64 `json:"offering_row"`
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
