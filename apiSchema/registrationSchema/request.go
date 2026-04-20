package registrationSchema

type RegisterStudentRequest struct {
	StudentID  int64 `json:"student_id"`
	OfferingID int64 `json:"offering_row"`
}
