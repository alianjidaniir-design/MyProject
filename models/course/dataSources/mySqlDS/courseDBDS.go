package mySqlDS

import (
	"MyProject/apiSchema/courseSchema"
	courseDataModle "MyProject/models/course/dataModels"
	courseDataSources "MyProject/models/course/dataSources"
	"MyProject/pkg/pagination"
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
func (ds *CourseDBDS) UpdateCourse(ctx context.Context, req courseSchema.UpdateCourseRequest) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	now := time.Now().In(myLocation())
	updateQuery := fmt.Sprintf("UPDATE %s SET updated_at = ? WHERE id = ?", ds.tableSQL)
	update, err := ds.db.PrepareContext(ctx, updateQuery)
	if err != nil {
		return course, err
	}
	defer update.Close()
	result, err := update.ExecContext(ctx, now, req.ID)
	if err != nil {
		return course, err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return courseDataModle.Course{}, err
	}
	var createdAt, updatedAt, deletedAt sql.NullTime
	readQuery := fmt.Sprintf("SELECT id , course_code , title , capacity ,enrolled_at ,isActive , created_at , updated_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)
	if err = ds.db.QueryRowContext(ctx, readQuery, req.ID).Scan(&course.ID, &course.CourseCode, &course.Title, &course.Capacity, &course.EnrolledAt, &course.IsActive, &createdAt, &updatedAt, &deletedAt); err != nil {
		return courseDataModle.Course{}, err
	}
	if createdAt.Valid {
		course.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		course.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		course.DeletedAt = deletedAt.Time
	}
	return course, nil
}

func (ds *CourseDBDS) GetCourse(ctx context.Context, req courseSchema.GetCoursesRequest) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	readQuery := fmt.Sprintf("SELECT id , course_code , title , capacity ,enrolled_at ,isActive , created_at , updated_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)
	var createdAt, updatedAt, deletedAt sql.NullTime
	if err := ds.db.QueryRowContext(ctx, readQuery, req.ID).Scan(&course.ID, &course.CourseCode, &course.Title, &course.Capacity, &course.EnrolledAt, &course.IsActive, &createdAt, &updatedAt, &deletedAt); err != nil {
		return courseDataModle.Course{}, err
	}
	if createdAt.Valid {
		course.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		course.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		course.DeletedAt = deletedAt.Time
	}
	return course, nil

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

func (ds *CourseDBDS) ListCourse(ctx context.Context, req courseSchema.CoursesListRequest) ([]courseDataModle.Course, int64, error) {
	var courses []courseDataModle.Course
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return courses, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var total int64
	totalItem := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableSQL)
	err = ds.db.QueryRowContext(ctx, totalItem).Scan(&total)
	if err != nil {
		return courses, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT id , course_code , title , capacity ,enrolled_at ,isActive ,  created_at, updated_at, deleted_at   FROM %s LIMIT ? OFFSET ?", ds.tableSQL)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return courses, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var course courseDataModle.Course
		var createdAt, updatedAt, deletedAt sql.NullTime
		err = rows.Scan(&course.ID, &course.CourseCode, &course.Title, &course.Capacity, &course.EnrolledAt, &course.IsActive, &createdAt, &updatedAt, &deletedAt)
		if err != nil {
			return courses, 0, err
		}
		if createdAt.Valid {
			course.CreatedAt = createdAt.Time.In(myLocation())
		}
		if updatedAt.Valid {
			course.UpdatedAt = updatedAt.Time.In(myLocation())
		}
		if deletedAt.Valid {
			course.DeletedAt = deletedAt.Time.In(myLocation())
		}
		courses = append(courses, course)

	}

	if rows.Err() != nil {
		return courses, 0, err

	}

	return courses, total, nil

}

func (ds *CourseDBDS) readCouresByID(ctx context.Context, id int64) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	readQuery := fmt.Sprintf("SELECT id , course_code , title , capacity ,enrolled_at ,isActive , created_at , updated_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)
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

func (ds *CourseDBDS) DeleteCourse(ctx context.Context, req courseSchema.HardDeleteCourseRequest) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id=?", ds.tableSQL)
	_, err := ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return course, err
	}
	return course, nil
}
