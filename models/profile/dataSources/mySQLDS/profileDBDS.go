package mySQLDS

import (
	"MyProject/apiSchema/profileSchema"
	"MyProject/models/profile/dataModels"
	"MyProject/pkg/pagination"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type ProfileDBDS struct {
	tableName string
	db        *sql.DB
}

func NewProfileDBDS(tableName string, db *sql.DB) (*ProfileDBDS, error) {
	cfg := &ProfileDBDS{
		tableName: tableName,
		db:        db,
	}
	return cfg, nil
}

func (ds *ProfileDBDS) CreateScoreStudent(ctx context.Context, req profileSchema.CreateScoresReq) (res dataModels.Profile, err error) {
	var profile dataModels.Profile
	tx, err := ds.db.BeginTx(ctx, nil)
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		} else if err != nil {
			tx.Rollback()
		}
	}()
	var checking bool
	checkRegisterQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM registration WHERE ID = ? AND deleted_at IS NULL AND status = ? ) THEN 1 ELSE 0 END
`
	var enroll = constants.StatusEnrolled
	err = tx.QueryRowContext(ctx, checkRegisterQuery, req.RegistrationID, enroll).Scan(&checking)
	if err != nil {
		return dataModels.Profile{}, err
	}
	if !checking {
		return dataModels.Profile{}, errors.New("Registration does not exist")
	}
	insertQuery := `
INSERT INTO profiles (registration_id , status_score , grade , score) VALUES (?, ?, ?, ?)
`
	switch {
	case req.Score > 20 || req.Score < 0:
		return dataModels.Profile{}, errors.New("Score must be between 0 and 20")
	case req.Score > 17:
		profile.Grade = "A"
		profile.StatusScore = constants.Passed
	case req.Score > 14:
		profile.Grade = "B"
		profile.StatusScore = constants.Passed
	case req.Score > 10:
		profile.Grade = "C"
		profile.StatusScore = constants.Passed
	case req.Score > 7:
		profile.Grade = "D"
		profile.StatusScore = constants.Failed
	case req.Score < 7:
		profile.Grade = "E"
		profile.StatusScore = constants.Failed
	default:
		return dataModels.Profile{}, errors.New("invalid score")
	}
	lastID, err := tx.ExecContext(ctx, insertQuery, req.RegistrationID, profile.StatusScore, profile.Grade, req.Score)
	if err != nil {
		return dataModels.Profile{}, err
	}

	result, err := lastID.LastInsertId()
	if err != nil {
		return dataModels.Profile{}, err
	}
	err = tx.Commit()
	if err != nil {
		return dataModels.Profile{}, err
	}
	return ds.readOProfileByID(ctx, result)
}

func (ds *ProfileDBDS) ListScoresStudents(crx context.Context, req profileSchema.ListAllScoresReq) (res []dataModels.ScoresStudents, total int, err error) {
	var profile []dataModels.ScoresStudents
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalRows int
	countQuery := `SELECT COUNT(*) FROM profiles`
	err = ds.db.QueryRowContext(crx, countQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, err
	}
	selectQuery := `
SELECT
u.ID    AS student_id,
u.code AS student_code,
c.ID AS course_id,
c.course_number AS course_number,
o.row AS offering_row,
o.group_number AS offering_group_number,
o.teacher_id AS offering_teacher_id,
p.status_score ,
p.grade ,
p.score 
FROM profiles p 
JOIN registration r ON p.registration_id = r.ID 
JOIN offerings o ON r.offering_row = o.row 
JOIN courses c ON o.course_id = c.ID 
JOIN student u ON r.student_id = u.ID 
ORDER BY u.code LIMIT ? OFFSET ?;
`
	rows, err := ds.db.QueryContext(crx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var s dataModels.ScoresStudents
		err = rows.Scan(
			&s.StudentID,
			&s.StudentCode,
			&s.CourseID,
			&s.CourseNumber,
			&s.OfferingRows,
			&s.OfferingGroup,
			&s.OfferingTeacher,
			&s.StatusScore,
			&s.Grade,
			&s.Score,
		)
		if err != nil {
			return nil, 0, err
		}
		profile = append(profile, s)

	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return profile, totalRows, nil
}

func (ds *ProfileDBDS) ListSummeryStudents(ctx context.Context, req profileSchema.ListAllScoresReq) (res []dataModels.StudentsSummary, total int, err error) {
	var profile []dataModels.StudentsSummary
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalRows int
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM profiles")
	err = ds.db.QueryRowContext(ctx, countQuery).Scan(&totalRows)
	if err != nil {
		return nil, 0, errors.New("error getting total rows")
	}
	selectQuery := `
SELECT
u.ID    AS student_id,
u.name AS student_name,
u.family  AS student_family,
u.major AS major
COUNT(DISTINCT c.course_number) AS total_course,
AVG(p.score) AS average_score,
 CASE
        WHEN AVG(p.score) >= 17 THEN 'A'
        WHEN AVG(p.score) >= 14 THEN 'B'
        WHEN AVG(p.score) >= 10 THEN 'C'
        WHEN AVG(p.score) >= 7 THEN 'D'
        ELSE 'E'
    END AS total_grade,
SUM(c.unit) AS total_units
FROM profiles p
JOIN registration r ON p.registration_id = r.ID
JOIN offerings o ON r.offering_row = o.row
JOIN courses c ON o.course_id = c.ID
JOIN student u ON r.student_id = u.ID
GROUP BY u.ID , u.name, u.family, u.major 
ORDER BY u.code LIMIT ? OFFSET ?;
`
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("error in syntax", err)
	}
	defer rows.Close()
	for rows.Next() {
		var s dataModels.StudentsSummary
		err = rows.Scan(
			&s.StudentID,
			&s.StudentName,
			&s.StudentFamily,
			&s.Major,
			&s.TotalCourse,
			&s.AverageScore,
			&s.TotalGrade,
			&s.TotalUnits,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("error", err)
		}
		profile = append(profile, s)
	}
	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}
	return profile, totalRows, nil

}

func (ds ProfileDBDS) GetStudent(ctx context.Context, req profileSchema.GetScoresReq) (res []dataModels.ScoresAnnouncement, err error) {
	var profile []dataModels.ScoresAnnouncement
	selectQuery := `
SELECT
p.ID AS profile_id,
u.code AS student_code,
u.name AS student_name,
u.family  AS student_family,
u.major AS major,
o.group_number AS offering_group_number,
c.course_number AS course_number,
c.title AS course_title,
c.unit AS unit,
t.name AS teacher_name,
t.last_name AS teacher_last_name,
p.status_score,
p.grade,
p.score,
0 AS total_units,
AVG(p.score) AS average_score,
 CASE
        WHEN AVG(p.score) >= 17 THEN 'A'
        WHEN AVG(p.score) >= 14 THEN 'B'
        WHEN AVG(p.score) >= 10 THEN 'C'
        WHEN AVG(p.score) >= 7 THEN 'D'
        ELSE 'E'
    END AS total_grade
FROM profiles p
JOIN registration r ON p.registration_id = r.ID
JOIN offerings o ON r.offering_row = o.row
JOIN courses c ON o.course_id = c.ID
JOIN teachers t ON o.teacher_id = t.ID
JOIN student u ON r.student_id = u.ID
WHERE u.ID = ?
GROUP BY
    p.ID, 
    u.code,
    u.name,
    u.family,
    u.major,
    o.group_number, 
    c.course_number,
    c.title,
    c.unit,
    t.name,
    t.last_name,
    p.status_score,
    p.grade,
    p.score,
    u.ID

UNION ALL 
SELECT
    0 AS profile_id,
    "" AS student_code,
    "" AS student_name,
    "" AS student_family,
    "" AS major,
    0 AS offering_group_number,
    0 AS course_number,
    "" AS course_title,
    0 AS unit,
    "" AS teacher_name,
    "" AS teacher_last_name,
    "" AS status_score,
    "" AS grade,
    0 AS score,
    COALESCE(SUM(c.unit), 0) AS total_units, -- اصلاح شده: استفاده از COALESCE
    0 AS average_score,
    0 AS total_grade
 

    FROM registration r
    JOIN offerings o ON r.offering_row = o.row
    JOIN courses c ON o.course_id = c.ID
    JOIN student u ON r.student_id = u.ID
    WHERE u.ID = ? 

;
`
	rows, err := ds.db.QueryContext(ctx, selectQuery, req.StudentID, req.StudentID)
	if err != nil {
		return res, fmt.Errorf("error in syntax", err)
	}
	defer rows.Close()
	for rows.Next() {
		var s dataModels.ScoresAnnouncement
		err = rows.Scan(&s.ID, &s.StudentCode, &s.StudentName, &s.StudentFamily, &s.Major, &s.OfferingGroupNumber, &s.CourseNumber, &s.CourseTitle, &s.Unit, &s.TeacherName, &s.TeacherLastName, &s.StatusScore, &s.Grade, &s.Score, &s.TotalUnits, &s.AverageScore, &s.TotalGrade)
		if err != nil {
			return res, fmt.Errorf("error in scanning", err)
		}
		profile = append(profile, s)

	}
	err = rows.Err()
	if err != nil {
		return res, fmt.Errorf("error in syntax", err)
	}
	return profile, nil
}

func (ds *ProfileDBDS) DeleteProfile(ctx context.Context, req profileSchema.DeleteScoresReq) (err error) {
	err = ds.checkID(ctx, req.ID)
	if err != nil {
		return err
	}
	deleted := fmt.Sprintf("DELETE FROM %s WHERE id=?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleted, req.ID)
	if err != nil {
		return fmt.Errorf("error in deletion", err)
	}
	return nil
}

func (ds *ProfileDBDS) readOProfileByID(ctx context.Context, ID int64) (res dataModels.Profile, err error) {
	var profile dataModels.Profile
	readQuery := fmt.Sprintf("SELECT ID , registration_id , status_score , grade  , score FROM %s WHERE ID = ? ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, readQuery, ID).Scan(&profile.ID, &profile.RegistrationID, &profile.StatusScore, &profile.Grade, &profile.Score)
	if err != nil {
		if err == sql.ErrNoRows {
			return dataModels.Profile{}, errors.New(sql.ErrNoRows.Error())
		}
		return dataModels.Profile{}, err
	}
	return profile, nil
}

func (ds *ProfileDBDS) checkID(ctx context.Context, row int64) error {
	var check bool
	searchQuery := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM profiles WHERE ID = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, searchQuery, row).Scan(&check)
	if err != nil {
		return errors.New(err.Error())
	}
	if !check {
		return errors.New("This is not a valid Row")
	}
	return nil
}
