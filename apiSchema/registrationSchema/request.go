package registrationSchema

type RegisterStudentRequest struct {
	StudentID  string `json:"student_id"`
	OfferingID string `json:"offering_id"`
}
