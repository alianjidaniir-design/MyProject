package mySqlDS

import (
	"MyProject/apiSchema/courseSchema"
	courseDataModle "MyProject/models/course/dataModels"
	courseDataSources "MyProject/models/course/dataSources"
	"context"
	"database/sql"
	"fmt"
	"time"
)

type CourseDBDS struct {
	tableName string
	tableSQL  string
	db        *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		fmt.Println(err)
	}
	return loc
}

func NewCourseDBDS(tableName string, db *sql.DB) (courseDataSources.CourseDB, error) {
	ff := &CourseDBDS{
		tableName: tableName,
		tableSQL:  tableName,
		db:        db,
	}
	return ff, nil
}

func (ds *CourseDBDS) CreateCourse(ctx context.Context, req courseSchema.RequestCourse) (courseDataModle.Course, error) {
	now := time.Now().In(myLocation())
	insertQuery := fmt.Sprintf("INSERT INTO %s (course_code, title , capacity , isActive , created_at , updated_at) VALUES (?, ?, ?, ? , ? , ?)", ds.tableSQL)
	insertresult, err := ds.db.ExecContext(ctx, insertQuery, req.CourseCode, req.Title, req.Capacity, req.IsActive, now, now)
	if err != nil {
		return courseDataModle.Course{}, err
	}
	insertID, err := insertresult.LastInsertId()
	if err != nil {
		return courseDataModle.Course{}, err
	}
	return ds.readCouresByID(ctx, insertID)

}

func (ds *CourseDBDS) readCouresByID(ctx context.Context, id int64) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	readQuery := fmt.Sprintf("SELECT id , course_code , title , capacity ,enrolled_count ,isActive , created_at , updated_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)
	var createdAt, updatedAt, deletedAt sql.NullTime
	if err := ds.db.QueryRowContext(ctx, readQuery, id).Scan(&course.ID, &course.CourseCode, &course.Title, &course.Capacity, &course.EnrolledAt, &course.IsActive, &createdAt, &updatedAt, &deletedAt); err != nil {
		return courseDataModle.Course{}, err
	}
	if createdAt.Valid {
		course.CreatedAt = createdAt.Time
	} else {
		course.CreatedAt = time.Time{}
	}
	if updatedAt.Valid {
		course.UpdatedAt = updatedAt.Time
	} else {
		course.UpdatedAt = time.Time{}
	}
	if deletedAt.Valid {
		course.DeletedAt = deletedAt.Time
	} else {
		course.DeletedAt = time.Time{}
	}
	return course, nil
}
