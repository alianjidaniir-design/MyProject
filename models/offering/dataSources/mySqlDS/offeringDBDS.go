package mySqlDS

import (
	"MyProject/apiSchema/offeringSchema"
	"MyProject/models/offering/dataModels"
	"MyProject/pkg/filter"
	"MyProject/pkg/timeLoc"
	"MyProject/pkg/val"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type OfferingDBDS struct {
	tableName string
	db        *sql.DB
}

func NewOfferingDBDS(tableName string, db *sql.DB) (*OfferingDBDS, error) {
	offer := &OfferingDBDS{
		tableName: tableName,
		db:        db,
	}
	return offer, nil
}

func (ds *OfferingDBDS) CreateOffering(ctx context.Context, req offeringSchema.CreateOfferingRequest) (res dataModels.Offering, err error) {
	err = val.CheckValidation(req)
	if err != nil {
		return dataModels.Offering{}, err
	}
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
CASE WHEN EXISTS (SELECT 1 FROM courses WHERE course_number = ? ) THEN 1 ELSE 0 END
`
	err = tx.QueryRowContext(ctx, courseQuery, req.CourseNumber).Scan(&checkCourse)
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
	timeStart := req.ClassStartTime
	timeEnd := req.ClassEndTime
	start, err := timeLoc.FormatTime(timeStart)
	if err != nil {
		return dataModels.Offering{}, err
	}
	end, err := timeLoc.FormatTime(timeEnd)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if start > end {
		return dataModels.Offering{}, errors.New("start time is greater than end time")
	}

	checkUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE group_number = ?", ds.tableName)
	err = tx.QueryRowContext(ctx, checkUnique, req.GroupNumber).Scan(&check)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if check > 0 {
		return dataModels.Offering{}, errors.New("groupNumber already exists")
	}

	var count int
	if req.Week == constants.Even {
		checkingUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE teacher_id = ? AND (week = ? OR week = ?) AND class_start_time < ? AND class_end_time > ?", ds.tableName)
		err = tx.QueryRowContext(ctx, checkingUnique, req.TeacherId, constants.Even, constants.All, start, end).Scan(&count)
		if err != nil {

			return dataModels.Offering{}, err
		}
		if count > 0 {
			return dataModels.Offering{}, errors.New("groupNumber already exists or It conflicts with this teacher's other class.")
		}
	} else if req.Week == constants.Odd {
		checkingUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE teacher_id = ? AND (week = ? OR week = ?) AND day = ? AND class_start_time <= ? AND class_end_time >= ?", ds.tableName)
		err = tx.QueryRowContext(ctx, checkingUnique, req.TeacherId, constants.Odd, constants.All, req.Day, start, end).Scan(&count)
		if err != nil {
			return dataModels.Offering{}, err
		}
		if count > 0 {
			return dataModels.Offering{}, errors.New("groupNumber already exists or It conflicts with this teacher's other class.")
		}
	} else if req.Week == constants.All {
		checkingUnique := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE teacher_id = ? AND day = ? AND class_start_time < ? AND class_end_time > ?", ds.tableName)
		err = tx.QueryRowContext(ctx, checkingUnique, req.TeacherId, req.Day, start, end).Scan(&count)
		if err != nil {
			return dataModels.Offering{}, err
		}
		if count > 0 {
			return dataModels.Offering{}, errors.New("groupNumber already exists or It conflicts with this teacher's other class.")
		}
	} else {
		return dataModels.Offering{}, errors.New("invalid week")
	}

	newID := lastID + 1
	insertQuery := fmt.Sprintf("INSERT INTO %s (row , group_number , course_number , teacher_id , capacity , isActive ,term_id,week, day , class_start_time , class_end_time, exam_start_time, exam_finish_time ) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?)", ds.tableName)
	_, err = tx.ExecContext(ctx, insertQuery, newID, req.GroupNumber, req.CourseNumber, req.TeacherId, req.Capacity, req.IsActive, req.TermId, req.Week, req.Day, start, end, req.ExamStartTime, req.ExamEndTime)
	if err != nil {
		return dataModels.Offering{}, err
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Offering{}, err
	}
	return ds.readOfferingByID(ctx, newID)

}

func (ds *OfferingDBDS) ListOffering(ctx context.Context, req offeringSchema.ListOfferingsRequest) (res []dataModels.ListOfferings, total int, err error) {
	var offerings []dataModels.ListOfferings
	var totalRows int
	err = val.CheckValidation(req)
	if err != nil {
		return nil, 0, err
	}
	fil := []filter.Filter{

		{Con: "college", Value: req.College},
		{Con: "educational_group", Value: req.EducationalGroup},
		{Con: "week", Value: req.Week},
		{Con: "day", Value: req.Day},
	}
	cond, args := filter.Filtering(fil...)
	whereClause := "term = ? AND year = ?"
	queryArgs := []interface{}{req.Term, req.Year}

	if cond != "" && cond != "1=1" {
		whereClause += " AND " + cond
		queryArgs = append(queryArgs, args...)
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM offering_list WHERE %s", whereClause)
	err = ds.db.QueryRowContext(ctx, countQuery, queryArgs...).Scan(&totalRows)
	if err != nil {
		return nil, 0, fmt.Errorf("error in rows count: %w", err)
	}
	selectQuery := fmt.Sprintf("SELECT * FROM offering_list WHERE %s ORDER BY year DESC, term DESC ", whereClause)
	rows, err := ds.db.QueryContext(ctx, selectQuery, queryArgs...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to fetch rows: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var offering dataModels.ListOfferings
		err = rows.Scan(&offering.Row, &offering.GroupNumber, &offering.CourseNumber, &offering.Title, &offering.Unit, &offering.TeacherName, &offering.TeacherLastName, &offering.College, &offering.EducationalGroup, &offering.Capacity, &offering.EnrolledCount, &offering.IsActive, &offering.Reservation, &offering.Week, &offering.Day, &offering.ClassStartTime, &offering.ClassEndTime, &offering.ExamStartTime, &offering.ExamEndTime, &offering.Term, &offering.Year)
		if err != nil {
			return nil, 0, fmt.Errorf("Error scanning row", err.Error())
		}
		offerings = append(offerings, offering)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, fmt.Errorf("Error : ", err.Error())
	}
	return offerings, totalRows, nil
}

func (ds *OfferingDBDS) GetOffering(ctx context.Context, req offeringSchema.GetRowOfferingRequest) (res dataModels.Offering, err error) {
	err = ds.checkID(ctx, req.Row)
	if err != nil {
		return res, err
	}
	return ds.readOfferingByID(ctx, req.Row)
}
func (ds *OfferingDBDS) DeActiveOffering(ctx context.Context, req offeringSchema.GetRowOfferingRequest) (res dataModels.Offering, err error) {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ? AND isActive = true ) THEN 1 ELSE 0 END
`
	err = ds.db.QueryRowContext(ctx, search, req.Row).Scan(&check)
	if err != nil {
		return dataModels.Offering{}, err
	}
	if !check {
		return dataModels.Offering{}, errors.New("active Offering does not exist")
	}
	deActiveQuery := fmt.Sprintf("UPDATE `%s` SET isActive = 0 WHERE row = ?", ds.tableName)
	update, err := ds.db.PrepareContext(ctx, deActiveQuery)
	if err != nil {
		return dataModels.Offering{}, err
	}
	defer update.Close()
	result, err := update.ExecContext(ctx, req.Row)
	if err != nil {
		return dataModels.Offering{}, err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows == 0 {
		return dataModels.Offering{}, err
	}
	return ds.readOfferingByID(ctx, req.Row)
}

func (ds *OfferingDBDS) readOfferingByID(ctx context.Context, row int64) (res dataModels.Offering, err error) {
	var offering dataModels.Offering
	readQuery := fmt.Sprintf("SELECT row , group_number , course_number , teacher_id , capacity , enrolled_count , isActive , reserveation , term_id , week , day , class_start_time , class_end_time , exam_start_time , exam_finish_time FROM %s WHERE row = ? ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, readQuery, row).Scan(&offering.Row, &offering.GroupNumber, &offering.CourseNumber, &offering.TeacherID, &offering.Capacity, &offering.EnrolledCount, &offering.IsActive, &offering.Reservation, &offering.TermID, &offering.Week, &offering.Day, &offering.ClassStartTime, &offering.ClassEndTime, &offering.ExamStartTime, &offering.ExamEndTime)
	if err != nil {
		if err == sql.ErrNoRows {
			return dataModels.Offering{}, errors.New(sql.ErrNoRows.Error())
		}
		return dataModels.Offering{}, err
	}
	return offering, nil
}

func (ds *OfferingDBDS) checkID(ctx context.Context, row int64) error {
	var check bool
	searchQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM offerings WHERE row = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, searchQuery, row).Scan(&check)
	if err != nil {
		return errors.New(err.Error())
	}
	if !check {
		return errors.New("This is not a valid ID")
	}
	return nil
}
