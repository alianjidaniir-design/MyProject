package mySqlDS

import (
	"MyProject/apiSchema/offeringSchema"
	"MyProject/models/offering/dataModels"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type OfferingDBDS struct {
	tableName string
	db        *sql.DB
}

func MyLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/Tehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}
func NewOfferingDBDS(tableName string, db *sql.DB) (*OfferingDBDS, error) {
	offer := &OfferingDBDS{
		tableName: tableName,
		db:        db,
	}
	return offer, nil
}

func (ds *OfferingDBDS) CreateOffering(ctx context.Context, req offeringSchema.CreateOfferingRequest) (res dataModels.Offering, err error) {
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return res, err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if err != nil {
			tx.Rollback()
		}
	}()
	var checkCourse bool
	courseQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM courses WHERE ID = ? ) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, courseQuery, req.CourseId).Scan(&checkCourse)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !checkCourse {
		return dataModels.Offering{}, errors.New("course does not exist")
	}

	var checkingTeacher bool

	teacherQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM teachers WHERE id = ?) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, teacherQuery, req.TeacherId).Scan(&checkingTeacher)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !checkingTeacher {
		return res, errors.New("Teacher does not exist")
	}

	var checkTerm bool
	termQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM terms WHERE id = ?) THEN 1 ELSE 0 END`
	err = tx.QueryRowContext(ctx, termQuery, req.TermId).Scan(&checkTerm)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !checkTerm {
		return res, errors.New("Terms does not exist")
	}

	var lastID int64

	lastIDQuery := fmt.Sprintf("SELECT COALESCE(MAX(row), 0) FROM %s", ds.tableName)
	err = tx.QueryRowContext(ctx, lastIDQuery).Scan(&lastID)
	if err != nil {
		return dataModels.Offering{}, err
	}
	var check int
	checkUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE group_number = ?", ds.tableName)
	err = tx.QueryRowContext(ctx, checkUnique, req.GroupNumber).Scan(&check)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if check > 0 {
		return dataModels.Offering{}, errors.New("groupNumber already exists")
	}

	newID := lastID + 1
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , group_number , course_id , teacher_id , capacity , isActive ,term_id, class_start_time , class_end_time, exam_start_time, exam_finish_time ) VALUES (?,?,?,?,?,?,?,?,?,?,?)", ds.tableName)
	_, err = tx.ExecContext(ctx, insertQuery, newID, req.GroupNumber, req.CourseId, req.TeacherId, req.Capacity, req.IsActive, req.TermId, req.ClassStartTime, req.ClassEndTime, req.ExamStartTime, req.ExamEndTime)
	if err != nil {
		tx.Rollback()
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Offering{}, err
	}
	return ds.readOfferingByID(ctx, newID)

}

func (ds *OfferingDBDS) readOfferingByID(ctx context.Context, row int64) (res dataModels.Offering, err error) {
	var offering dataModels.Offering
	readQuery := fmt.Sprintf("SELECT row , group_number , course_id , teacher_id , capacity , enrolled_count , isActive , reserveation , term_id , class_start_time , class_end_time , exam_start_time , exam_finish_time FROM %s WHERE row = ? ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, readQuery, row).Scan(&offering.Row, &offering.GroupNumber, &offering.CourseID, &offering.TeacherID, &offering.Capacity, &offering.EnrolledCount, &offering.IsActive, &offering.Reservation, &offering.TermID, &offering.ClassStartTime, &offering.ClassEndTime, &offering.ExamStartTime, &offering.ExamEndTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return dataModels.Offering{}, errors.New(sql.ErrNoRows.Error())
		}
		return dataModels.Offering{}, err
	}
	return offering, nil
}
