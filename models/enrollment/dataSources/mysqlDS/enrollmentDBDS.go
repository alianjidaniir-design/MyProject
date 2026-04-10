package mysqlDS

import (
	"MyProject/apiSchema/commonSchema"
	"MyProject/apiSchema/enrollmentSchema"
	courseDataModle "MyProject/models/course/dataModels"
	EnrollmentDataModel "MyProject/models/enrollment/dataModels"
	"MyProject/models/enrollment/dataSources"
	UserdataModel "MyProject/models/user/dataModel"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type EnrollmentDBDS struct {
	tableName string
	db        *sql.DB
}

func myLocation() *time.Location {
	location, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return location
}

func NewEnrollmentDBDS(tablename string, db *sql.DB) (dataSources.EnrollmentDS, error) {
	ff := &EnrollmentDBDS{
		tableName: tablename,
		db:        db,
	}
	return ff, nil

}

func (ds *EnrollmentDBDS) EnrollStudent(ctx context.Context, req enrollmentSchema.EnrollmentRequest) (res EnrollmentDataModel.Enrollment, err error) {
	now := time.Now().In(myLocation())
	var enrollment EnrollmentDataModel.Enrollment
	var course courseDataModle.Course
	var student UserdataModel.User

	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}
	defer tx.Rollback()
	var ExitStudent, ExirCourse bool
	queryCombinedLogic := `
    SELECT 
        CASE WHEN EXISTS(SELECT 1 FROM student WHERE student_id = ?) THEN TRUE ELSE FALSE END AS student_ok,
        CASE WHEN EXISTS(SELECT 1 FROM courses WHERE course_id = ? ) THEN TRUE ELSE FALSE END AS course_ok
`
	err = tx.QueryRowContext(ctx, queryCombinedLogic, req.StudentID, req.CourseID).Scan(&ExitStudent, &ExirCourse)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	} else if !ExitStudent {
		return EnrollmentDataModel.Enrollment{}, errors.New("student does not exist")
	} else if !ExirCourse {
		return EnrollmentDataModel.Enrollment{}, errors.New("course does not exist")
	}
	var checkActiviate bool
	querycombin := `
SELECT
CASE WHEN EXISTS(SELECT 1 FROM courses WHERE course_id = ? AND isActive = true) THEN 1 ELSE 0 END AS check_activiate
`
	err = tx.QueryRowContext(ctx, querycombin, req.CourseID).Scan(&checkActiviate)
	if err != nil {
		if !checkActiviate {
			return EnrollmentDataModel.Enrollment{}, errors.New("this course is deActive")
		}
		return EnrollmentDataModel.Enrollment{}, err
	}
	var enrollmentOK bool
	queryenroll := `
SELECT
CASE WHEN EXISTS(SELECT 1 FROM enrollments WHERE id = ?) THEN TRUE ELSE 0 END AS enrollment_ok
`
	err = tx.QueryRowContext(ctx, queryenroll, req.CourseID).Scan(&enrollmentOK)
	if err != nil {
		if enrollmentOK {
			return EnrollmentDataModel.Enrollment{}, errors.New("student alReady enrolled ")
		}
		return EnrollmentDataModel.Enrollment{}, err
	}
	var checkingcapacity bool
	checkcapacity := `
SELECT
CASE WHEN EXISTS(SELECT 1 FROM courses WHERE course_id = ? AND capacity > enrolled_at) THEN TRUE ELSE FALSE AS capacity_ok
`
	err = tx.QueryRowContext(ctx, checkcapacity, req.CourseID).Scan(&checkingcapacity)
	if err != nil {
		if !checkingcapacity {
			return EnrollmentDataModel.Enrollment{}, errors.New("this class 's capacity is full")
		}
		return EnrollmentDataModel.Enrollment{}, err
	}

	sqlStatement := fmt.Sprintf("INSERT INTO %s (student_id, course_id, enrolledd_at, created_at, updated_at, deleted_at) VALUES (?, ?, ?, ?, ?, ?)", ds.tableName)
	_, err = tx.ExecContext(ctx, sqlStatement, req.StudentID, req.CourseID, now, now, now, nil)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}
	updateenrollment := fmt.Sprintf("UPDATE %s SET enrolledd_at = enrolledd_at + 1 WHERE id = ?", ds.tableName)
	_, err = tx.ExecContext(ctx, updateenrollment, enrollment.ID)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}

}
