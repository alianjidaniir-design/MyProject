package dataModels

type Profile struct {
	ID             int64  `json:"row"`
	RegistrationID int64  `json:"registration_id"`
	StatusScore    string `json:"status_score"`
	Grade          string `json:"grade"`
	Score          int    `json:"score"`
}

type ScoresStudents struct {
	StudentID       int64  `json:"student_id"`
	StudentCode     string `json:"student_code"`
	CourseID        int64  `json:"course_id"`
	CourseNumber    string `json:"course_number"`
	OfferingRows    int64  `json:"offering_rows"`
	OfferingGroup   int    `json:"offering_group"`
	OfferingTeacher int64  `json:"offering_teacher"`
	StatusScore     string `json:"status_score"`
	Grade           string `json:"grade"`
	Score           int    `json:"score"`
}

type StudentsSummary struct {
	StudentID     int64   `json:"student_id"`
	StudentName   string  `json:"student_name"`
	StudentFamily string  `json:"student_family"`
	Major         string  `json:"major"`
	TotalCourse   int     `json:"total_course"`
	AverageScore  float64 `json:"average_score"`
	TotalUnits    int     `json:"total_units"`
	TotalGrade    string  `json:"total_grade"`
}

type 
