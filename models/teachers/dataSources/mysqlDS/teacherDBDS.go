package mysqlDS

import (
	"MyProject/apiSchema/teacherSchema"
	"MyProject/pkg/hash"
	"MyProject/pkg/pagination"
	TimeLoc "MyProject/pkg/timeLoc"
	Val "MyProject/pkg/val"
	"MyProject/statics/constants/roles"
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
	err = Val.CheckValidation(req)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	hashing, err := hash.HashingPassword(req.Password)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	roleName := roles.RoleTeacher
	var count int
	now := time.Now().In(TimeLoc.MyLocation())
	query := `SELECT COUNT(*) FROM student WHERE national_code = ? `
	err = ds.db.QueryRowContext(ctx, query, req.NationalCode).Scan(&count)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	if count > 0 {
		return dataModels.Teacher{}, errors.New("national code exists")
	}
	query = `SELECT COUNT(*) FROM teachers WHERE national_code = ? `
	err = ds.db.QueryRowContext(ctx, query, req.NationalCode).Scan(&count)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	if count > 0 {
		return dataModels.Teacher{}, errors.New("national code exists")
	}
	insert := fmt.Sprintf("INSERT INTO %s (name , last_name , role_name , national_code , email , phone , work_experience , password, created_at , updated_at ) VALUES (?, ?, ?, ?, ?, ?, ? , ? , ? , ?)", ds.tableName)
	insertQuery, err := ds.db.ExecContext(ctx, insert, req.Name, req.LastName, roleName, req.NationalCode, req.Email, req.Phone, req.WorkExperience, hashing, now, now)
	if err != nil {
		return dataModels.Teacher{}, err
	}

	insertID, err := insertQuery.LastInsertId()
	if err != nil {
		return dataModels.Teacher{}, err
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
	selectQuery := fmt.Sprintf("SELECT ID , name , last_name , email , phone , work_experience , created_at , updated_at , deleted_at FROM %s  LIMIT ? OFFSET ? ", ds.tableName)
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
	err = ds.chackTeacher(ctx, req.ID)
	if err != nil {
		return "", err
	}
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ? ", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return res, err
	}
	response := "deleted done successfully"
	return response, nil
}

func (ds *TeacherDBDS) SoftDeleteTeachers(ctx context.Context, req teacherSchema.SelectTeacherSchema) (res dataModels.Teacher, err error) {
	err = ds.chackTeacher(ctx, req.ID)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	now := time.Now().In(TimeLoc.MyLocation())
	updateQuery := fmt.Sprintf("UPDATE %s SET deleted_at=? , updated_at = ? WHERE id=?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, updateQuery, now, now, req.ID)
	if err != nil {
		return dataModels.Teacher{}, err
	}

	return ds.readQuery(ctx, req.ID)
}
func (ds *TeacherDBDS) UpdateTeachers(ctx context.Context, req teacherSchema.SelectTeacherSchema) (res dataModels.Teacher, err error) {
	now := time.Now().In(TimeLoc.MyLocation())
	err = ds.chackTeacher(ctx, req.ID)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	update := fmt.Sprintf("UPDATE %s SET updated_at = ? WHERE id = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, update, now, req.ID)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	return ds.readQuery(ctx, req.ID)

}

func (ds *TeacherDBDS) LoginTeachers(ctx context.Context, req teacherSchema.LoginTeacherRequest) (res string, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return "", err
	}
}

func (ds *TeacherDBDS) chackTeacher(ctx context.Context, ID int64) error {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM teachers WHERE ID = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, search, ID).Scan(&check)

	if err != nil {
		return err
	}
	if !check {
		return errors.New("Teacher not found")
	}
	return nil
}

func (ds *TeacherDBDS) readQuery(ctx context.Context, ID int64) (dataModels.Teacher, error) {
	var teacher dataModels.Teacher
	read := fmt.Sprintf("SELECT ID , name , last_name , role_name , national_code , email , phone , work_experience , password , created_at , updated_at , deleted_at  FROM %s WHERE ID=? ", ds.tableName)

	err := ds.db.QueryRowContext(ctx, read, ID).Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.RoleName, &teacher.NationalCode, &teacher.Email, &teacher.Phone, &teacher.WorkExperience, &teacher.Password, &teacher.CreatedAt, &teacher.UpdatedAt, &teacher.DeletedAt)
	if err != nil {
		return teacher, err
	}

	return teacher, nil
}

func (ds *TeacherDBDS) checkingTeacher(ctx context.Context, email string, code string) (dataModels.Teacher, error) {
	var teacher dataModels.Teacher
	checkQuery := "SELECT * FROM teachers WHERE email = ? AND national_code = ?"
	err := ds.db.QueryRowContext(ctx, checkQuery, email, code).Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.RoleName, &teacher.NationalCode, &teacher.Email, &teacher.Phone, &teacher.WorkExperience, &teacher.Password, &teacher.CreatedAt, &teacher.UpdatedAt, &teacher.DeletedAt)
	if err != nil {
		return teacher, err
	}
	return teacher, nil
}
