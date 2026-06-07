package registrationSchema

import "MyProject/models/Registrations/dataModels"

type RegisterStudentResponse struct {
	Information []dataModels.ListSelectOfferingResponse
}

type DetailStudentResponse struct {
	Information dataModels.Registration
}

type ClassSchedule struct {
	MyClasses []dataModels.TermClassSchedules
	Page      int `json:"page"`
	Total     int
}

type ClassesTeacher struct {
	MyClasses []dataModels.TermClasses
	Page      int `json:"page"`
	Total     int
}
type DeleteStudentResponse struct {
	Information dataModels.Registration
	Massage     string
}

type ListStudentsResponse struct {
	List  []dataModels.Registration
	Total int
}

type ListStudentResponse struct {
	List  []dataModels.Offering
	Total int
}

type ListOfferingResponse struct {
	List  []dataModels.Student
	Total int
}

type CancelRegistrationResponse struct {
	Information dataModels.Registration
	Massage     string
}
