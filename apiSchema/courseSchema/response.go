package courseSchema

import (
	courseDataModle "MyProject/models/course/dataModels"
)

type ResponseCourse struct {
	Course courseDataModle.Course `json:"course"`
}
