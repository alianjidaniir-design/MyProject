package enrollmentSchema

import "MyProject/models/enrollment/dataModels"

type EnrollmentResponse struct {
	Enrollment dataModels.Enrollment `json:"enrollment"`
}
