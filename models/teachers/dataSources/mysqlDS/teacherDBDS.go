package mysqlDS

import (
	"MyProject/apiSchema/teacherSchema"
	tokenDataModel "MyProject/models/token/dataModel"
	"MyProject/pkg/hash"
	"MyProject/pkg/pagination"
	TimeLoc "MyProject/pkg/timeLoc"
	"MyProject/pkg/token"
	Val "MyProject/pkg/val"
	"MyProject/statics/constants"
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
	var check bool
	checkRole := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM roles WHERE name = ?) THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, checkRole, req.RoleName).Scan(&check)
	if err != nil {
		return dataModels.Teacher{}, err
	}
	if !check {
		return dataModels.Teacher{}, errors.New("role does not exist")
	}
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
	if req.Password != "0"+req.NationalCode {
		return dataModels.Teacher{}, errors.New("password and national code do not match")
	}
	if count > 0 {
		return dataModels.Teacher{}, errors.New("national code exists")
	}
	insert := fmt.Sprintf("INSERT INTO %s (name , last_name , role_name , national_code , email , phone , work_experience , password, created_at , updated_at ) VALUES (?, ?, ?, ?, ?, ?, ? , ? , ? , ?)", ds.tableName)
	insertQuery, err := ds.db.ExecContext(ctx, insert, req.Name, req.LastName, req.RoleName, req.NationalCode, req.Email, req.Phone, req.WorkExperience, hashing, now, now)
	if err != nil {
		return dataModels.Teacher{}, err
	}

	insertID, err := insertQuery.LastInsertId()
	if err != nil {
		return dataModels.Teacher{}, err
	}
	return ds.readQuery(ctx, insertID)
}

func (ds *TeacherDBDS) MyInfo(ctx context.Context, ID int64) (dataModels.InfoTeacher, error) {
	err := ds.chackTeacher(ctx, ID)
	if err != nil {
		return dataModels.InfoTeacher{}, err
	}
	var info dataModels.InfoTeacher
	selectQuery := fmt.Sprintf("SELECT name , last_name ,  national_code , email , phone ,  work_experience FROM %s WHERE ID = ?", ds.tableName)
	err = ds.db.QueryRowContext(ctx, selectQuery, ID).Scan(&info.Name, &info.LastName, &info.NationalCode, &info.Email, &info.Phone, &info.WorkExperience)
	if err != nil {
		return dataModels.InfoTeacher{}, err
	}
	return info, nil
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
	selectQuery := fmt.Sprintf("SELECT ID , name , last_name ,role_name,national_code, email , phone , work_experience ,password, created_at , updated_at , deleted_at FROM %s  LIMIT ? OFFSET ? ", ds.tableName)
	rows, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return []dataModels.Teacher{}, 0, errors.New("there is an error in pagination")
	}
	defer rows.Close()
	for rows.Next() {
		var teacher dataModels.Teacher
		var createdAt, updatedAt, deletedAt sql.NullTime
		err = rows.Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.RoleName, &teacher.NationalCode, &teacher.Email, &teacher.Phone, &teacher.WorkExperience, &teacher.Password, &createdAt, &updatedAt, &deletedAt)
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

func (ds *TeacherDBDS) LoginTeachers(ctx context.Context, req teacherSchema.LoginTeacherRequest) (access string, refresh string, massage string, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return "", "", "", err
	}
	student, err := ds.checkingTeacher(ctx, req.Email, req.NationalCode)
	if err != nil {
		return "", "", "", err
	}
	err = hash.CheckPassword(req.Password, student.Password)
	if err != nil {
		return "", "", "", err
	}
	tokenAccess, err := token.GenerateAccessToken(student.ID, student.RoleName)
	if err != nil {
		return "", "", "", err
	}
	RefreshToken, err := token.GenerateRefreshToken(student.RoleName, student.ID)
	if err != nil {
		return "", "", "", err
	}
	expires := time.Now().In(TimeLoc.MyLocation()).Add(constants.RefreshTokenExpiry)
	insertQuery := fmt.Sprintf("INSERT INTO refreshs (user_id , role_name , token , expires_at , rekoved_at) VALUES (?, ?, ?, ?, ?)")
	_, err = ds.db.QueryContext(ctx, insertQuery, student.ID, student.RoleName, RefreshToken, expires, false)
	if err != nil {
		return "", "", "", err
	}
	return tokenAccess, RefreshToken, fmt.Sprintf("Welcome %s ", student.Name), nil
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
	checkQuery := "SELECT ID , name , last_name , role_name , national_code , email , phone , work_experience , password , created_at , updated_at , deleted_at  FROM teachers WHERE email = ? AND national_code = ?"
	err := ds.db.QueryRowContext(ctx, checkQuery, email, code).Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.RoleName, &teacher.NationalCode, &teacher.Email, &teacher.Phone, &teacher.WorkExperience, &teacher.Password, &teacher.CreatedAt, &teacher.UpdatedAt, &teacher.DeletedAt)
	if err != nil {
		return teacher, err
	}
	return teacher, nil
}

func (ds *TeacherDBDS) Refresh(ctx context.Context, req string) (access string, refresh string, err error) {
	var rt tokenDataModel.RefreshToken
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return "", "", err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
	}()
	if req == "" {
		return "", "", errors.New("refresh Token Required")
	}
	checkToken := "SELECT user_id, role_name, token, expires_at, rekoved_at FROM refreshs WHERE token = ? AND rekoved_at = false"
	err = tx.QueryRowContext(ctx, checkToken, req).Scan(&rt.UserID, &rt.RoleName, &rt.Token, &rt.ExpiresAt, &rt.RevokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == true {
			return "", "", errors.New("Refresh token is invalid")
		}
		return "", "", err
	}
	currentTime := time.Now().In(TimeLoc.MyLocation())
	deleteQuery := "DELETE FROM refreshs WHERE expires_at < ?"
	_, err = tx.ExecContext(ctx, deleteQuery, currentTime)
	if err != nil {
		return "", "", err
	}

	student, err := ds.readQuery(ctx, rt.UserID)
	if err != nil {
		return "", "", err
	}
	newToken, err := token.GenerateRefreshToken(student.RoleName, student.ID)
	if err != nil {
		return "", "", err
	}
	updated := fmt.Sprintf("UPDATE refreshs SET rekoved_at = true WHERE token = ?")
	rows, err := tx.PrepareContext(ctx, updated)
	if err != nil {
		return "", "", err
	}
	defer rows.Close()
	_, err = rows.ExecContext(ctx, rt.Token)
	if err != nil {
		return "", "", err
	}

	expiresAt := time.Now().In(TimeLoc.MyLocation()).Add(constants.RefreshTokenExpiry)

	insertQuery := fmt.Sprintf("INSERT INTO refreshs (user_id , role_name , token , expires_at , rekoved_at) VALUES (?, ?, ? , ? , ?)")
	result, err := tx.ExecContext(ctx, insertQuery, rt.UserID, rt.RoleName, newToken, expiresAt, false)
	if err != nil {
		return "", "", err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return "", "", errors.New("failed to revoke old token")
	}

	accessToken, err := token.GenerateAccessToken(student.ID, student.RoleName)
	if err != nil {
		return "", "", err
	}
	if err = tx.Commit(); err != nil {
		return "", "", err
	}
	return accessToken, newToken, nil
}

func (ds *TeacherDBDS) Logout(ctx context.Context, req teacherSchema.LogoutSchema, refresh string) error {
	var rt tokenDataModel.RefreshToken
	tx, err := ds.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			return
		}
		if err != nil {
			tx.Rollback()
			return
		}
	}()
	if req.IsLogout == false {
		return errors.New("request is canceled")
	}
	if refresh == "" {
		return errors.New("token is empty")
	}
	checkToken := fmt.Sprintf("SELECT user_id, role_name, token, expires_at, rekoved_at FROM refreshs WHERE token = ? AND rekoved_at = false")
	err = tx.QueryRowContext(ctx, checkToken, refresh).Scan(&rt.UserID, &rt.RoleName, &rt.Token, &rt.ExpiresAt, &rt.RevokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == true {
			return errors.New("Refresh token is invalid")
		}
		return err

	} else if time.Now().In(TimeLoc.MyLocation()).After(rt.ExpiresAt) {
		return errors.New("token expired")
	}

	rt.RevokedAt = true
	updated := fmt.Sprintf("UPDATE refreshs SET rekoved_at = ? WHERE token = ?")
	rows, err := tx.PrepareContext(ctx, updated)
	if err != nil {
		return err
	}
	defer rows.Close()
	row, err := rows.ExecContext(ctx, rt.RevokedAt, rt.Token)
	if err != nil {
		return err
	}
	rowsAffected, err := row.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated, token may not exist")
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
