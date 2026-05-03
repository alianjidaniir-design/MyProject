package tuitionSchema

import "MyProject/models/tuition/dataModels"

type InformationTuitionSchema struct {
	Detail dataModels.Tuition
}

type MassageTuition struct {
	Detail  dataModels.Tuition
	Massage string
}

type ListTuitionSchema struct {
	Detail []dataModels.StudentsDebit
	Total  int
}

type ListAllTuitionSchema struct {
	Detail []dataModels.Tuition
	Total  int
}

type TuitionStudentSchema struct {
	Detail            []dataModels.TuitionStudent
	TotalUnits        int
	TotalDebitAmount  int
	TotalCreditAmount int
	Reminder          int
}
