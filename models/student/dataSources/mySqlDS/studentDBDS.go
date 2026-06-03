package mySqlDS

import (
	"MyProject/apiSchema/studentSchema"
	studentDataModel "MyProject/models/student/dataModel"
	studentDataSources "MyProject/models/student/dataSources"
	tokenDataModel "MyProject/models/token/dataModel"
	"MyProject/pkg/hash"
	"MyProject/pkg/pagination"
	TimeLoc "MyProject/pkg/timeLoc"
	"MyProject/pkg/token"
	Val "MyProject/pkg/val"
	"MyProject/statics/constants"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type StudentDBDS struct {
	tableName string
	tableSQL  string
	db        *sql.DB
}

func NewStudentDBDS(db *sql.DB, tableName string) (studentDataSources.StudentDB, error) {

	userDBInstance := &StudentDBDS{
		tableName: tableName,
		tableSQL:  tableName,
		db:        db,
	}
	return userDBInstance, nil
}

func (ds *StudentDBDS) SoftDeleteStudent(ctx context.Context, req studentSchema.SoftDeleteRequest) (studentDataModel.Student, error) {
	err := Val.CheckValidation(req)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	now := time.Now().In(TimeLoc.MyLocation())
	err = ds.chackStudent(ctx, req.ID)
	if err != nil {
		return studentDataModel.Student{}, errors.New(err.Error())
	}
	updateQuery := fmt.Sprintf("UPDATE %s SET updated_at=? , deleted_at=? WHERE id = ?", ds.tableName)

	_, err = ds.db.ExecContext(ctx, updateQuery, now, now, req.ID)
	if err != nil {
		return studentDataModel.Student{}, err
	}

	return ds.readTaskByID(ctx, req.ID)
}

func (ds *StudentDBDS) DeleteStudent(ctx context.Context, req studentSchema.DeleteRequest) (studentDataModel.Student, error) {
	err := ds.chackStudent(ctx, req.ID)
	if err != nil {
		return studentDataModel.Student{}, errors.New("Found Not student")
	}
	var students studentDataModel.Student
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", ds.tableName)
	_, err = ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return studentDataModel.Student{}, err
	}

	return students, nil
}

func (ds *StudentDBDS) GetStudent(ctx context.Context, req studentSchema.GetRequest) (studentDataModel.Student, error) {
	err := ds.chackStudent(ctx, req.ID)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	return ds.readTaskByID(ctx, req.ID)
}

func (ds *StudentDBDS) CreateStudent(ctx context.Context, req studentSchema.SignUpStudent) (studentDataModel.Student, error) {
	err := Val.CheckValidation(req)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	ha, err := hash.HashingPassword(req.Password)
	code := &req.StudentCode
	req.UserName = code
	var check, check2 bool
	checkRole := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM roles WHERE name = ?) THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, checkRole, req.RoleName).Scan(&check)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	if !check {
		return studentDataModel.Student{}, errors.New("Invalid role")
	}
	checkTerm := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM terms WHERE ID = ?) THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, checkTerm, req.TermID).Scan(&check2)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	if !check {
		return studentDataModel.Student{}, errors.New("Invalid term")
	}
	now := time.Now().In(TimeLoc.MyLocation())
	insertQuery := fmt.Sprintf("INSERT INTO %s (name , family, phone ,national_code, major,student_code, term_id , level ,user_name ,password, role_name , created_at , updated_at , deleted_at) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", ds.tableSQL)
	insertResult, err := ds.db.ExecContext(ctx, insertQuery, req.Name, req.Family, req.Phone, req.NationalCode, req.Major, req.StudentCode, req.TermID, req.Level, code, ha, req.RoleName, now, now, nil)
	if err != nil {
		return studentDataModel.Student{}, err
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return studentDataModel.Student{}, err
	}
	return ds.readTaskByID(ctx, insertedID)
}

func (ds *StudentDBDS) ReadStudent(ctx context.Context, req studentSchema.ListRequest) ([]studentDataModel.Student, int64, error) {
	var users []studentDataModel.Student // نام متغیر به جمع تغییر یافت
	Page, PageSize, err := pagination.CheckPage(req.Page, req.PageSize)
	if err != nil {
		return nil, 0, err
	}
	// "offest" به "offset" اصلاح شد
	offset := (Page - 1) * PageSize
	limit := PageSize
	var total int64
	totalItem := fmt.Sprintf("SELECT COUNT(*) FROM %s", ds.tableSQL)
	err = ds.db.QueryRowContext(ctx, totalItem).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// ستون‌ها را صریحاً نام ببرید تا از مشکلات احتمالی ترتیب ستون‌ها جلوگیری شود.
	// فرض می‌کنیم ترتیب ستون‌ها در دیتابیس با ترتیب مدل مطابقت دارد.
	selectQuery := fmt.Sprintf("SELECT id, name, family,phone,national_code,major,student_code,term_id , level,user_name,password,role_name, created_at, updated_at, deleted_at FROM %s LIMIT ? OFFSET ?", ds.tableSQL)
	selectResult, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return []studentDataModel.Student{}, 0, err
	}
	defer selectResult.Close()

	for selectResult.Next() {
		var student studentDataModel.Student
		// تعریف متغیرهای موقت از نوع sql.NullTime برای دریافت مقادیر NULL پذیر

		// اسکن مقادیر از دیتابیس به متغیرهای موقت NullTime
		if err = selectResult.Scan(&student.ID, &student.Name, &student.Family, &student.Phone, &student.NationalCode, &student.Major, &student.StudentCode, &student.UserName, &student.Password, &student.RoleName, &student.CreatedAt, &student.UpdatedAt, &student.DeletedAt); err != nil {
			// اگر اینجا خطا رخ داد، ممکن است به دلیل عدم تطابق نوع یا نام ستون باشد
			return nil, 0, fmt.Errorf("خطا در اسکن ردیف: %w", err)
		}

		users = append(users, student) // اضافه کردن به slice 'users'
	}
	if err = selectResult.Err(); err != nil {
		return nil, 0, fmt.Errorf("خطا در پیمایش نتایج کوئری: %w", err)
	}
	return users, total, nil
}

func (ds *StudentDBDS) RenameStudent(ctx context.Context, req studentSchema.UpdateUserRequest) (studentDataModel.Student, error) {
	var students studentDataModel.Student
	stmt := fmt.Sprintf("UPDATE %s SET  name = ?, family = ?, updated_at = ? WHERE id = ? ", ds.tableName)
	var updatedAt time.Time
	sss, err := ds.db.PrepareContext(ctx, stmt)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	defer sss.Close()

	result, err := sss.ExecContext(ctx,
		students.Name,
		students.Family,
		updatedAt,
		req.ID,
	)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	// (optional) require for number of updated column
	rows, err := result.RowsAffected()
	if err != nil {
		return studentDataModel.Student{}, errors.New("error in number update")
	}
	if rows == 0 {
		return studentDataModel.Student{}, fmt.Errorf("rows == 0")
	}
	updatedAt = updatedAt.In(TimeLoc.MyLocation())
	return students, nil

}

func (ds *StudentDBDS) UpdateStudent(ctx context.Context, req studentSchema.UpdateUserRequest) (studentDataModel.Student, error) {
	now := time.Now().In(TimeLoc.MyLocation())
	err := ds.chackStudent(ctx, req.ID)
	if err != nil {
		return studentDataModel.Student{}, errors.New("Found Not student")
	}
	stmt := fmt.Sprintf("UPDATE %s SET updated_at = ? WHERE id = ? ", ds.tableName)
	sss, err := ds.db.PrepareContext(ctx, stmt)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	defer sss.Close()

	result, err := sss.ExecContext(ctx,
		now,
		req.ID,
	)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	// (optional) require for number of updated column
	rows, err := result.RowsAffected()
	if err != nil {
		return studentDataModel.Student{}, errors.New("error in number update")
	}
	if rows == 0 {
		return studentDataModel.Student{}, fmt.Errorf("rows == 0")
	}

	return ds.readTaskByID(ctx, req.ID)

}

func (ds *StudentDBDS) StudentEntry(ctx context.Context, req studentSchema.LoginStudent) (string, string, string, error) {
	err := Val.CheckValidation(req)
	if err != nil {
		return "", "", "", err
	}
	student, err := ds.checkingStudent(req.UserName)
	if err != nil {

		return "", "", "", err
	}
	err = hash.CheckPassword(req.Password, student.Password)
	if err != nil {
		fmt.Println(req.Password, student.Password, 12)

		return "", "", "", err
	}
	tok, err := token.GenerateAccessToken(student.ID, student.RoleName)
	if err != nil {
		fmt.Println("s4")

		return "", "", "", err
	}
	refresh, err := token.GenerateRefreshToken(student.RoleName, student.ID)
	if err != nil {
		return "", "", "", err
	}

	expiresAt := time.Now().In(TimeLoc.MyLocation()).Add(constants.RefreshTokenExpiry)

	insert := fmt.Sprintf("INSERT INTO refreshs (user_id,role_name , token , expires_at , rekoved_at ) VALUES (?, ?, ? ,?, ?)")
	if _, err = ds.db.ExecContext(ctx, insert, student.ID, student.RoleName, refresh, expiresAt, false); err != nil {
		return "", "", "", err
	}
	massage := fmt.Sprintf("Welcome %s", student.Name)
	return tok, refresh, massage, nil

}

func (ds *StudentDBDS) RefreshToken(ctx context.Context, req string) (string, string, error) {
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
		return "", "", errors.New("refresh token cannot be empty")
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

	student, err := ds.readTaskByID(ctx, rt.UserID)
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
	_, err = tx.ExecContext(ctx, insertQuery, rt.UserID, rt.RoleName, newToken, expiresAt, false)
	if err != nil {
		return "", "", err
	}

	accessToken, err := token.GenerateAccessToken(student.ID, student.RoleName)
	if err != nil {
		return "", "", err
	}
	if err = tx.Commit(); err != nil {
		return "", "", err
	}
	return newToken, accessToken, nil

}

func (ds *StudentDBDS) readTaskByID(ctx context.Context, userID int64) (studentDataModel.Student, error) {
	var students studentDataModel.Student

	readQuery := fmt.Sprintf("SELECT id , name , family,phone,national_code,major,student_code,term_id,level,user_name , password,role_name , created_at , updated_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)

	if err := ds.db.QueryRowContext(ctx, readQuery, userID).Scan(&students.ID, &students.Name, &students.Family, &students.Phone, &students.NationalCode, &students.Major, &students.StudentCode, &students.TermID, &students.Level, &students.UserName, &students.Password, &students.RoleName, &students.CreatedAt, &students.UpdatedAt, &students.DeletedAt); err != nil {
		return studentDataModel.Student{}, err
	}

	return students, nil

}
func (ds *StudentDBDS) chackStudent(ctx context.Context, ID int64) error {
	var check bool
	search := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM student WHERE ID = ?) THEN 1 ELSE 0 END
`
	err := ds.db.QueryRowContext(ctx, search, ID).Scan(&check)

	if err != nil {
		return err
	}
	if !check {
		return errors.New("Student not found")
	}
	return nil
}

func (ds *StudentDBDS) checkingStudent(s string) (data studentDataModel.Student, err error) {

	var students studentDataModel.Student
	selectQuery := fmt.Sprintf("SELECT ID , name , student_code , password , role_name FROM student WHERE student_code = ?")
	err = ds.db.QueryRow(selectQuery, s).Scan(&students.ID, &students.Name, &students.StudentCode, &students.Password, &students.RoleName)
	if err != nil {
		return studentDataModel.Student{}, err
	}
	return students, nil
}

func (ds *StudentDBDS) RevokedRefreshToken(ctx context.Context, req studentSchema.LogoutRequest, tok string) error {
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
	if tok == "" {
		return errors.New("token is empty")
	}
	checkToken := fmt.Sprintf("SELECT user_id, role_name, token, expires_at, rekoved_at FROM refreshs WHERE token = ? AND rekoved_at = false")
	err = tx.QueryRowContext(ctx, checkToken, tok).Scan(&rt.UserID, &rt.RoleName, &rt.Token, &rt.ExpiresAt, &rt.RevokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == true {
			fmt.Println(tok)
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
