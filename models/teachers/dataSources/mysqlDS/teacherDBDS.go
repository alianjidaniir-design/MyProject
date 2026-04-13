package mysqlDS

import (
	"MyProject/apiSchema/teacherSchema"
	"MyProject/pkg/pagination"
	"fmt"

	"MyProject/models/teachers/dataModels"
	"MyProject/models/teachers/dataSources"
	"context"
	"database/sql"
	"errors"
	"time"
)

type TeacherDBDS struct {
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

func NewTeacherDBDS(tableName string, db *sql.DB) (dataSources.TeacherDS, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	ff := &TeacherDBDS{
		tableName: tableName,
		db:        db,
	}

	return ff, nil

}

func (ds *TeacherDBDS) CreateTeacher(ctx context.Context, req teacherSchema.InformationSchema) (res dataModels.Teacher, err error) {
	var teacher dataModels.Teacher
	now := time.Now().In(myLocation())
	checkQueryEmail := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE email=? ", ds.tableName)
	var emailCount, phoneCount int
	err = ds.db.QueryRowContext(ctx, checkQueryEmail, req.Email).Scan(&emailCount)
	if err != nil || err == sql.ErrNoRows {
		return dataModels.Teacher{}, fmt.Errorf("there is a error getting the email and phone count ", err)
	}
	if emailCount > 0 {
		return dataModels.Teacher{}, fmt.Errorf("there are email ")
	}
	checkQueryPhone := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE phone=? ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, checkQueryPhone, req.Phone).Scan(&phoneCount)
	if err != nil || err == sql.ErrNoRows {
		return dataModels.Teacher{}, fmt.Errorf("there is a error getting the phone count ", err)
	}
	if phoneCount > 0 {
		return dataModels.Teacher{}, fmt.Errorf("there are phone")
	}
	insert := fmt.Sprintf("INSERT INTO %s (name , last_name , email , phone , work_experience , created_at , updated_at ) VALUES (?, ?, ?, ?, ?, ?, ?)", ds.tableName)

	insertQuery, err := ds.db.ExecContext(ctx, insert, req.Name, req.LastName, req.Email, req.Phone, req.Description, now, now)
	if err != nil {
		return teacher, err
	}
	insertID, err := insertQuery.LastInsertId()
	if err != nil {
		return teacher, err
	}
	return ds.readQuery(ctx, insertID)

}

func (ds *TeacherDBDS) ListTeachers(ctx context.Context, req teacherSchema.PaginationSchema) (res []dataModels.Teacher, total int64, err error) {
	var teachers []dataModels.Teacher
	page, pageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return []dataModels.Teacher{}, 0, errors.New("there is an error checking the page and page size")
	}
	offset := (page - 1) * pageSize
	limit := pageSize
	var totalAll int64
	totaling := fmt.Sprintf("SELECT COUNT(*) FROM %s ", ds.tableName)
	err = ds.db.QueryRowContext(ctx, totaling).Scan(&totalAll)
	if err != nil {
		return []dataModels.Teacher{}, 0, errors.New("there is an error in total the page and page size")
	}
	selectQuery := fmt.Sprintf("SELECT ID , name , last_name , email , phone , work_experience , created_at , updated_at , deleted_at FROM %s LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return []dataModels.Teacher{}, 0, errors.New("there is an error in pagination")
	}
	defer rows.Close()
	for rows.Next() {
		var teacher dataModels.Teacher
		var createdAt, updatedAt, deletedAt sql.NullTime
		err = rows.Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.Email, &teacher.Phone, &teacher.WorkExperience, &createdAt, &updatedAt, &deletedAt)
		if err != nil {
			return []dataModels.Teacher{}, 0, errors.New("there is an error for scanning the rows")
		}
		if createdAt.Valid {
			teacher.CreatedAt = createdAt.Time
		}
		if updatedAt.Valid {
			teacher.UpdatedAt = updatedAt.Time
		}
		teachers = append(teachers, teacher)

	}
	err = rows.Err()
	if err != nil {
		return []dataModels.Teacher{}, 0, err
	}
	return teachers, totalAll, nil

}
func (ds *TeacherDBDS) GetTeacherById(ctx context.Context, req teacherSchema.GetTeacherSchema) (res dataModels.Teacher, err error) {
	return ds.readQuery(ctx, req.ID)
}

func (ds *TeacherDBDS) HardDeleteTeachers(ctx context.Context, req teacherSchema.SelectTeacherSchema) (res string, err error) {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM teachers WHERE id=?) THEN 1 ELSE 0 END
`
	err = ds.db.QueryRowContext(ctx, search, req.ID).Scan(&check)
	if err != nil {
		return "", err
	}
	if !check {
		return "Teacher not found", nil
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ? ", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return res, err
	}
	response := "deleted done successfully"
	return response, nil
}

func (ds *TeacherDBDS) readQuery(ctx context.Context, ID int64) (dataModels.Teacher, error) {
	var teacher dataModels.Teacher
	read := fmt.Sprintf("SELECT ID , name , last_name , email , phone , work_experience , created_at , updated_at , deleted_at  FROM %s WHERE ID=?", ds.tableName)

	var createdAt, updatedAt, deletedAt sql.NullTime
	err := ds.db.QueryRowContext(ctx, read, ID).Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.Email, &teacher.Phone, &teacher.WorkExperience, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return teacher, err
	}
	if createdAt.Valid {
		teacher.CreatedAt = createdAt.Time
	}
	if updatedAt.Valid {
		teacher.UpdatedAt = updatedAt.Time
	}

	return teacher, nil
}
