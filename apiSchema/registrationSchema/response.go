package registrationSchema

import "MyProject/models/Registrations/dataModels"

type RegisterStudentResponse struct {
	Information dataModels.Registration
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

type CancelRegistrationResponse struct {
	Information dataModels.Registration
	Massage     string
}
