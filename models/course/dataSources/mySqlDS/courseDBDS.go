package mySqlDS

import (
	"MyProject/apiSchema/courseSchema"
	courseDataModle "MyProject/models/course/dataModels"
	courseDataSources "MyProject/models/course/dataSources"
	"MyProject/pkg/pagination"
	TimeLoc "MyProject/pkg/timeLoc"
	Val "MyProject/pkg/val"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type CourseDBDS struct {
	tableSQL string
	db       *sql.DB
}

func NewCourseDBDS(tableName string, db *sql.DB) (courseDataSources.CourseDB, error) {
	ff := &CourseDBDS{
		tableSQL: tableName,
		db:       db,
	}
	return ff, nil
}
func (ds *CourseDBDS) CreateCourse(ctx context.Context, req courseSchema.RequestCourse) (courseDataModle.Course, error) {
	now := time.Now().In(TimeLoc.MyLocation())
	err := Val.CheckValidation(req)
	if err != nil {
		return courseDataModle.Course{}, err
	}
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM departments WHERE ID = ?) THEN 1 ELSE 0 END
`
	err = ds.db.QueryRowContext(ctx, search, req.DepartmentID).Scan(&check)

	if err != nil {
		return courseDataModle.Course{}, err
	}
	if !check {
		return courseDataModle.Course{}, errors.New("Department not found")
	}

	insertQuery := fmt.Sprintf("INSERT INTO %s (course_number, title, unit ,department_id , prerequisite , necessary , created_at , updated_at ) VALUES (?, ?, ?, ?, ?, ?, ? , ?)", ds.tableSQL)
	lastID, err := ds.db.ExecContext(ctx, insertQuery, req.CourseNumber, req.Title, req.Unit, req.DepartmentID, req.Prerequisite, req.Necessary, now, now)

	if err != nil {
		return courseDataModle.Course{}, fmt.Errorf("there are a problem in top query", err)
	}

	newID, err := lastID.LastInsertId()
	if err != nil {
		return courseDataModle.Course{}, fmt.Errorf("there are a problem in top query", err)
	}

	return ds.readCourseByID(ctx, newID)

}
func (ds *CourseDBDS) UpdateCourse(ctx context.Context, req courseSchema.UpdateCourseRequest) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	now := time.Now().In(TimeLoc.MyLocation())
	err := ds.chackCourse(ctx, req.ID)
	if err != nil {
		return courseDataModle.Course{}, errors.New("there is not course")
	}
	updateQuery := fmt.Sprintf("UPDATE %s SET updated_at = ? WHERE ID = ?", ds.tableSQL)
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

	return ds.readCourseByID(ctx, req.ID)
}

func (ds *CourseDBDS) GetCourse(ctx context.Context, req courseSchema.GetCoursesRequest) (courseDataModle.Course, error) {
	err := ds.chackCourse(ctx, req.ID)
	if err != nil {
		return courseDataModle.Course{}, errors.New("Course not found")
	}
	return ds.readCourseByID(ctx, req.ID)

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
		return nil, 0, err
	}
	selectQuery := fmt.Sprintf("SELECT ID ,course_number , title , unit , department_id , prerequisite , necessary ,  created_at, updated_at, deleted_at FROM %s LIMIT ? OFFSET ?", ds.tableSQL)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var course courseDataModle.Course
		if err = rows.Scan(&course.ID, &course.CourseNumber, &course.Title, &course.Unit, &course.DepartmentID, &course.Prerequisite, &course.Necessary, &course.CreatedAt, &course.UpdatedAt, &course.DeletedAt); err != nil {
			return nil, 0, err
		}

		courses = append(courses, course)

	}

	if rows.Err() != nil {
		return []courseDataModle.Course{}, 0, err

	}

	return courses, total, nil

}

func (ds *CourseDBDS) ListDepartmentsCourse(ctx context.Context, req courseSchema.DepartmentListRequest) ([]courseDataModle.Course, int64, error) {
	var courses []courseDataModle.Course
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, errors.New("there is a error in checkPage")
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalPage int64
	totalItem := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE department_id = ?", ds.tableSQL)
	err = ds.db.QueryRowContext(ctx, totalItem, req.DepartmentID).Scan(&totalPage)
	if err != nil {
		return nil, 0, fmt.Errorf("there is a error in Query ", err)
	}
	if totalPage == 0 {
		return nil, 0, errors.New("there is no department or this is a a problem in counter")
	}
	selectQuery := fmt.Sprintf("SELECT ID ,course_number , title , unit , department_id , prerequisite , necessary ,  created_at, updated_at, deleted_at FROM %s WHERE department_id = ? LIMIT ? OFFSET ?", ds.tableSQL)
	rows, err := ds.db.QueryContext(ctx, selectQuery, req.DepartmentID, limit, offset)
	if err != nil {
		return nil, 0, errors.New("the pagination query failed")
	}
	defer rows.Close()
	for rows.Next() {
		var course courseDataModle.Course
		if err = rows.Scan(&course.CourseNumber, &course.Title, &course.Unit, &course.DepartmentID, &course.Prerequisite, &course.Necessary, &course.CreatedAt, &course.UpdatedAt, &course.DeletedAt); err != nil {
			return nil, 0, err
		}

		courses = append(courses, course)
	}
	if rows.Err() != nil {
		return nil, 0, err
	}
	return courses, totalPage, nil
}

func (ds *CourseDBDS) readCourseByID(ctx context.Context, id int64) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	readQuery := fmt.Sprintf("SELECT ID , course_number , title , unit , department_id , prerequisite , necessary, created_at , updated_at , deleted_at FROM %s WHERE ID = ? ORDER BY ID ", ds.tableSQL)
	if err := ds.db.QueryRowContext(ctx, readQuery, id).Scan(&course.ID, &course.CourseNumber, &course.Title, &course.Unit, &course.DepartmentID, &course.Prerequisite, &course.Necessary, &course.CreatedAt, &course.UpdatedAt, &course.DeletedAt); err != nil {
		return courseDataModle.Course{}, err
	}

	return course, nil
}

func (ds *CourseDBDS) DeleteCourse(ctx context.Context, req courseSchema.HardDeleteCourseRequest) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	err := ds.chackCourse(ctx, req.ID)
	if err != nil {
		return course, errors.New("Course Found not")
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE ID = ?", ds.tableSQL)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return course, err
	}
	return course, nil
}

func (ds *CourseDBDS) SoftDelete(ctx context.Context, req courseSchema.SoftDeleteCourseRequest) (courseDataModle.Course, error) {
	var course courseDataModle.Course
	now := time.Now().In(TimeLoc.MyLocation())
	err := ds.chackCourse(ctx, req.ID)
	if err != nil {
		return courseDataModle.Course{}, errors.New("Course Not Found")
	}
	update := fmt.Sprintf("UPDATE %s SET deleted_at = ? WHERE ID = ?", ds.tableSQL)
	_, err = ds.db.ExecContext(ctx, update, now, req.ID)
	if err != nil {
		return course, err
	}
	return ds.readCourseByID(ctx, req.ID)
}

func (ds *CourseDBDS) chackCourse(ctx context.Context, ID int64) error {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM courses WHERE ID = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, search, ID).Scan(&check)

	if err != nil {
		return err
	}
	if !check {
		return errors.New("Course not found")
	}
	return nil
}
