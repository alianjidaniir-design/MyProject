package dataModel

import "time"

type Activity struct {
	ActivityType string    `json:"activityType"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	TeacherID    int64     `json:"teacherID"`
	CourseID     int64     `json:"CourseIDd"`
	GroupNumber  int       `json:"groupNumber"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}
