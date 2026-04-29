package tuitionSchema

import "MyProject/models/tuition/dataModels"

type InformationTuitionSchema struct {
	Detail dataModels.Tuition
}

type MassageUpdateTuition struct {
	Detail  dataModels.Tuition
	Massage string
}
