package dataModels

type Profile struct {
	ID             int64  `json:"row"`
	RegistrationID int64  `json:"registration_id"`
	StatusScore    string `json:"status_score"`
	Grade          string `json:"grade"`
	Score          int    `json:"score"`
}

type ScoresResponse struct {
	StudentID       int64  `json:"student_id"`
	StudentCode     int    `json:"student_code"`
	CourseID        int64  `json:"course_id"`
	CourseNumber    int    `json:"course_number"`
	OfferingRows    int64  `json:"offering_rows"`
	OfferingGroup   int64  `json:"offering_group"`
	OfferingTeacher int64  `json:"offering_teacher"`
	StatusScore     string `json:"status_score"`
	Grade           string `json:"grade"`
	Score           int    `json:"score"`
}
