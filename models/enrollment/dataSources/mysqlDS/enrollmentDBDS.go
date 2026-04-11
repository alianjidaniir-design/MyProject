package mysqlDS

import (
	"MyProject/apiSchema/enrollmentSchema"
	EnrollmentDataModel "MyProject/models/enrollment/dataModels"
	"MyProject/models/enrollment/dataSources"
	"MyProject/pkg/pagination"
	"MyProject/statics/constants"
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
	if db == nil {
		return nil, errors.New("db is nil")
	}
	ff := &EnrollmentDBDS{
		tableName: tablename,
		db:        db,
	}

	return ff, nil

}

func (ds *EnrollmentDBDS) EnrollStudent(ctx context.Context, req enrollmentSchema.EnrollmentRequest) (res EnrollmentDataModel.Enrollment, err error) {

	now := time.Now().In(myLocation())
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}
	txCommitted := false

	defer func() {
		if !txCommitted {
			if r := recover(); r != nil {
				tx.Rollback()
				panic(r)
			} else if err != nil {
				tx.Rollback()
			}
		}
	}()
	var ExitStudent, ExirCourse bool
	queryCombinedLogic := `
    SELECT 
        CASE WHEN EXISTS(SELECT 1 FROM student WHERE id = ?) THEN TRUE ELSE FALSE END AS student_ok,
        CASE WHEN EXISTS(SELECT 1 FROM courses WHERE id = ? ) THEN TRUE ELSE FALSE END AS course_ok
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
	queryIsActive := `
        SELECT CASE WHEN isActive = true THEN 1 ELSE 0 END
        FROM courses
        WHERE id = ?
    `
	err = tx.QueryRowContext(ctx, queryIsActive, req.CourseID).Scan(&checkActiviate)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	} else if !checkActiviate {
		return EnrollmentDataModel.Enrollment{}, errors.New("this course is deActive")
	}
	var enrollmentOK bool
	queryenroll := `
SELECT
CASE WHEN EXISTS(SELECT 1 FROM enrollments WHERE student_id = ? AND course_id = ? AND canceled_at IS NULL) THEN TRUE ELSE FALSE END AS enrollment_ok
`
	err = tx.QueryRowContext(ctx, queryenroll, req.StudentID, req.CourseID).Scan(&enrollmentOK)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}
	if enrollmentOK {
		return EnrollmentDataModel.Enrollment{}, errors.New("student alReady enrolled ")
	}

	var checkingcapacity bool
	checkcapacity := `
SELECT
CASE WHEN EXISTS(SELECT 1 FROM courses WHERE id = ? AND capacity > enrolled_at) THEN 1 ELSE 0 END 

`
	err = tx.QueryRowContext(ctx, checkcapacity, req.CourseID).Scan(&checkingcapacity)
	if err != nil {

		return EnrollmentDataModel.Enrollment{}, err
	}

	if !checkingcapacity {
		return EnrollmentDataModel.Enrollment{}, errors.New("this class 's capacity is full")
	}

	sqlStatement := fmt.Sprintf("INSERT INTO %s (student_id, course_id,status, enrolledd_at, created_at, updated_at, deleted_at) VALUES (?,?, ?, ?, ?, ?, ?)", ds.tableName)
	var ff = constants.StatusEnrolled
	add, err := tx.ExecContext(ctx, sqlStatement, req.StudentID, req.CourseID, ff, now, now, now, nil)
	if err != nil {
		fmt.Println(sqlStatement)

		return EnrollmentDataModel.Enrollment{}, err
	}
	insert, err := add.LastInsertId()
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}
	updatenrolledCount := fmt.Sprintf("UPDATE courses SET enrolled_at = enrolled_at + 1 , updated_at = ? WHERE id = ?")
	_, err = tx.ExecContext(ctx, updatenrolledCount, now, req.CourseID)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}

	err = tx.Commit()
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}
	return ds.readQuery(ctx, insert)

}

func (ds *EnrollmentDBDS) CancelEnrollment(ctx context.Context, req enrollmentSchema.CancelEnrollmentRequest) (res EnrollmentDataModel.Enrollment, err error, result string) {
	now := time.Now().In(myLocation())
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, errors.New("error starting transaction"), ""
	}
	commitTx := false
	defer func() {
		if !commitTx {
			if e := recover(); e != nil {
				tx.Rollback()
				panic(e)
			} else if err != nil {
				tx.Rollback()
			}
		}
	}()
	var enrollmentOk bool
	var enrollee = constants.StatusEnrolled
	queryenroll := `
SELECT
CASE WHEN EXISTS(SELECT 1 FROM enrollments WHERE id = ? AND status = ?) THEN 1 ELSE 0 END 
`
	err = ds.db.QueryRowContext(ctx, queryenroll, req.ID, enrollee).Scan(&enrollmentOk)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err, ""
	}
	if !enrollmentOk {
		return EnrollmentDataModel.Enrollment{}, errors.New("student has not enrolled"), ""
	}
	var cancel = constants.StatusCanceled
	update := fmt.Sprintf("UPDATE %s SET canceled_at = ? , status = ? WHERE id = ? ", ds.tableName)
	_, err = ds.db.ExecContext(ctx, update, now, cancel, req.ID)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err, ""
	}
	var courseID int64
	query := fmt.Sprintf("SELECT course_id FROM enrollments WHERE id = ? FOR UPDATE ")
	err = tx.QueryRowContext(ctx, query, req.ID).Scan(&courseID)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err, ""
	}
	if courseID == 0 {
		tx.Rollback()
	}
	decrement := fmt.Sprintf("UPDATE courses SET enrolled_at = enrolled_at - 1 , updated_at = ? WHERE id = ? ")
	_, err = tx.ExecContext(ctx, decrement, now, courseID)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, errors.New("error starting transaction"), ""
	}

	err = tx.Commit()
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err, ""
	}
	var enrollment EnrollmentDataModel.Enrollment
	readQuery := fmt.Sprintf(`
        SELECT id, student_id, course_id, status, enrolledd_at, canceled_at, created_at, updated_at, deleted_at
        FROM %s
        WHERE id = ? AND status = ? `, ds.tableName)
	var ff = constants.StatusCanceled
	err = ds.db.QueryRowContext(ctx, readQuery, req.ID, ff).Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.Status, &enrollment.EnrolledAt, &enrollment.CanceledAt, &enrollment.CreatedAt, &enrollment.UpdatedAt, &enrollment.DeletedAt)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err, ""
	}

	return enrollment, nil, "stundet candeles successfully"

}
func (ds *EnrollmentDBDS) ListEnrollment(ctx context.Context, req enrollmentSchema.ListEnrollmentsRequest) (res []EnrollmentDataModel.Enrollment, code int64, err error) {
	var enroll []EnrollmentDataModel.Enrollment
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return []EnrollmentDataModel.Enrollment{}, 400, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var total int64
	countItem := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableName)
	err = ds.db.QueryRowContext(ctx, countItem).Scan(&total)
	if err != nil {
		return []EnrollmentDataModel.Enrollment{}, 400, errors.New("error getting enrollment count")
	}
	selectQuery := fmt.Sprintf("SELECT * FROM %s LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return []EnrollmentDataModel.Enrollment{}, 400, errors.New("error getting enrollment list")
	}
	defer rows.Close()
	for rows.Next() {
		var enrollment EnrollmentDataModel.Enrollment
		var createdAt, updatedAt, deletedAt sql.NullTime
		err = rows.Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.Status, &enrollment.EnrolledAt, &enrollment.CanceledAt, &createdAt, &updatedAt, &deletedAt)
		if err != nil {
			return []EnrollmentDataModel.Enrollment{}, 400, errors.New("error scaning enrollment list")
		}
		if createdAt.Valid {
			enrollment.CreatedAt = createdAt.Time.In(myLocation())
		} else {
			enrollment.CreatedAt = time.Time{}
		}
		if updatedAt.Valid {
			enrollment.UpdatedAt = updatedAt.Time.In(myLocation())
		} else {
			enrollment.UpdatedAt = time.Time{}
		}
		if deletedAt.Valid {
			deletedAt.Time = deletedAt.Time.In(myLocation())
		} else {
			deletedAt.Time = time.Time{}
		}
		enroll = append(enroll, enrollment)

	}
	if err = rows.Err(); err != nil {
		return []EnrollmentDataModel.Enrollment{}, 400, err
	}
	return enroll, total, nil
}

func (ds *EnrollmentDBDS) ListStudentCourses(ctx context.Context, req enrollmentSchema.ListStudentCoursesRequest) (res EnrollmentDataModel.Enrollment, err error) {
	var enrollment EnrollmentDataModel.Enrollment
	var courseID int64

	switch req.Status {
	case constants.StatusEnrolled:
		var enroll = constants.StatusEnrolled
		enrolled := fmt.Sprintf("SELECT course_id FROM %s WHERE status = ? AND student_id = ?", ds.tableName)
		_, err = ds.db.QueryContext(ctx, enrolled, enroll, req.StudentID)
		if err != nil {
			return EnrollmentDataModel.Enrollment{}, errors.New("error getting enrollment course")
		}
	case constants.StatusCanceled:
		var cancel = constants.StatusCanceled
		canceled := fmt.Sprintf("SELECT course_id FROM %s WHERE status = ? AND student_id = ?", ds.tableName)
		_, err = ds.db.QueryContext(ctx, canceled, cancel, req.StudentID)
		if err != nil {
			return EnrollmentDataModel.Enrollment{}, errors.New("error getting enrollment course")
		}
	case "":
		all := fmt.Sprintf("SELECT course_id FROM %s WHERE student_id = ?", ds.tableName)
		err = ds.db.QueryRowContext(ctx, all, req.StudentID).Scan(&courseID)
		if err != nil {
			return EnrollmentDataModel.Enrollment{}, errors.New("error getting enrollment course")
		}
	default:
		return EnrollmentDataModel.Enrollment{}, errors.New("status not supported")
	}
	return enrollment, nil
}

func (ds *EnrollmentDBDS) readQuery(ctx context.Context, ID int64) (EnrollmentDataModel.Enrollment, error) {
	var enrollment EnrollmentDataModel.Enrollment
	readQuery := fmt.Sprintf(`
        SELECT id, student_id, course_id, status, enrolledd_at, canceled_at, created_at, updated_at, deleted_at
        FROM %s
        WHERE id = ? AND status = ? `, ds.tableName)
	var ff = constants.StatusEnrolled
	err := ds.db.QueryRowContext(ctx, readQuery, ID, ff).Scan(&enrollment.ID, &enrollment.StudentID, &enrollment.CourseID, &enrollment.Status, &enrollment.EnrolledAt, &enrollment.CanceledAt, &enrollment.CreatedAt, &enrollment.UpdatedAt, &enrollment.DeletedAt)
	if err != nil {
		return EnrollmentDataModel.Enrollment{}, err
	}

	return enrollment, nil

}
