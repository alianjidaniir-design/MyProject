package mySqlDS

import (
	"MyProject/apiSchema/userSchema"
	userDataModel "MyProject/models/user/dataModel"
	userDataSourses "MyProject/models/user/dataSourses"
	"MyProject/pkg/pagination"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type UserDBDS struct {
	tableName string
	tableSQL  string
	db        *sql.DB
}

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/ُTehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func NewUsersDBDS(db *sql.DB, tableName string) (userDataSourses.UserDB, error) {

	userDBInstance := &UserDBDS{
		tableName: tableName,
		tableSQL:  tableName,
		db:        db,
	}
	return userDBInstance, nil
}

func (ds *UserDBDS) SoftDeleteStudent(ctx context.Context, req userSchema.SoftDeleteRequest) (userDataModel.User, error) {
	var student userDataModel.User
	now := time.Now().In(myLocation())
	updateQuery := fmt.Sprintf("UPDATE %s SET deleted_at=? WHERE id = ?", ds.tableName)

	_, err := ds.db.ExecContext(ctx, updateQuery, now, req.ID)
	if err != nil {
		return userDataModel.User{}, err
	}

	selectQuery := fmt.Sprintf("SELECT id, code, name, family, created_at, updated_at, deleted_at FROM %s WHERE id = ?", ds.tableName)

	var createdAtAl, updatedAtAl, deletedAtAl sql.NullTime

	err = ds.db.QueryRowContext(ctx, selectQuery, req.ID).Scan(
		&student.ID, &student.Code, &student.Name, &student.Family,
		&createdAtAl, &updatedAtAl, &deletedAtAl,
	)
	if err != nil {
		return userDataModel.User{}, fmt.Errorf("failed to fetch deleted user: %w", err)
	}

	if createdAtAl.Valid {
		student.CreatedAt = createdAtAl.Time.In(myLocation())
	}
	if updatedAtAl.Valid {
		student.UpdatedAt = updatedAtAl.Time.In(myLocation())
	}
	if deletedAtAl.Valid {
		student.DeletedAt = deletedAtAl.Time.In(myLocation())
	}

	return student, nil
}

func (ds *UserDBDS) DeleteStudent(ctx context.Context, req userSchema.DeleteRequest) (userDataModel.User, error) {

	var students userDataModel.User
	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = ?", ds.tableName)
	_, err := ds.db.ExecContext(ctx, deleteQuery, req.ID)
	if err != nil {
		return userDataModel.User{}, err
	}

	return students, nil
}

func (ds *UserDBDS) GetStudent(ctx context.Context, req userSchema.GetRequest) (userDataModel.User, error) {
	var students userDataModel.User
	selectQuery := fmt.Sprintf("SELECT * FROM %s WHERE id = ? ", ds.tableName)

	var createdAtAl, updatedAtAl, deletedAtAl sql.NullTime

	if err := ds.db.QueryRowContext(ctx, selectQuery, req.ID).Scan(&students.ID, &students.Code, &students.Name, &students.Family, &createdAtAl, &deletedAtAl, &updatedAtAl); err != nil {
		return userDataModel.User{}, err
	}

	if createdAtAl.Valid {
		students.CreatedAt = createdAtAl.Time.In(myLocation())
	} else {
		students.CreatedAt = time.Time{}
	}

	if updatedAtAl.Valid {
		students.UpdatedAt = updatedAtAl.Time.In(myLocation())
	} else {
		fmt.Println(updatedAtAl.Time.In(myLocation()))
		students.UpdatedAt = time.Time{}
	}
	if deletedAtAl.Valid {
		fmt.Println(deletedAtAl.Time.In(myLocation()))
		students.DeletedAt = deletedAtAl.Time.In(myLocation())
	} else {
		students.DeletedAt = time.Time{}
	}

	return students, nil

}

func (ds *UserDBDS) CreateStudent(ctx context.Context, req userSchema.LoginRequest) (userDataModel.User, error) {
	now := time.Now().In(myLocation())
	insertQuery := fmt.Sprintf("INSERT INTO %s (code , name , family , created_at , updated_at) VALUES (?, ? , ?,?,?)", ds.tableSQL)
	insertResult, err := ds.db.ExecContext(ctx, insertQuery, req.Code, req.Name, req.Family, now, now)
	if err != nil {
		return userDataModel.User{}, err
	}

	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return userDataModel.User{}, err
	}
	return ds.readTaskByID(ctx, insertedID)
}

func (ds *UserDBDS) ReadStudent(ctx context.Context, req userSchema.ListRequest) ([]userDataModel.User, int64, error) {
	var users []userDataModel.User // نام متغیر به جمع تغییر یافت
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
		return []userDataModel.User{}, 0, err
	}

	// ستون‌ها را صریحاً نام ببرید تا از مشکلات احتمالی ترتیب ستون‌ها جلوگیری شود.
	// فرض می‌کنیم ترتیب ستون‌ها در دیتابیس با ترتیب مدل مطابقت دارد.
	selectQuery := fmt.Sprintf("SELECT id, code, name, family, created_at, updated_at, deleted_at FROM %s LIMIT ? OFFSET ?", ds.tableSQL)
	selectResult, err := ds.db.QueryContext(ctx, selectQuery, limit, offset)
	if err != nil {
		return []userDataModel.User{}, 0, err
	}
	defer selectResult.Close()

	for selectResult.Next() {
		var user userDataModel.User
		// تعریف متغیرهای موقت از نوع sql.NullTime برای دریافت مقادیر NULL پذیر
		var createdAtSQL, updatedAtSQL, deletedAtSQL sql.NullTime

		// اسکن مقادیر از دیتابیس به متغیرهای موقت NullTime
		if err = selectResult.Scan(&user.ID, &user.Code, &user.Name, &user.Family, &createdAtSQL, &updatedAtSQL, &deletedAtSQL); err != nil {
			// اگر اینجا خطا رخ داد، ممکن است به دلیل عدم تطابق نوع یا نام ستون باشد
			return []userDataModel.User{}, 0, fmt.Errorf("خطا در اسکن ردیف: %w", err)
		}

		// تبدیل مقادیر NullTime به time.Time در مدل user
		// اگر مقدار معتبر (NULL نباشد) بود، آن را به فیلد مربوطه اختصاص دهید.
		// فرض می‌شود myLocation() تابع درستی برای دریافت منطقه زمانی است.
		if createdAtSQL.Valid {
			user.CreatedAt = createdAtSQL.Time.In(myLocation())
		} else {
			user.CreatedAt = time.Time{} // مقدار زمان صفر برای NULL
		}

		if updatedAtSQL.Valid {
			user.UpdatedAt = updatedAtSQL.Time.In(myLocation())
		} else {
			user.UpdatedAt = time.Time{} // مقدار زمان صفر برای NULL
		}

		if deletedAtSQL.Valid {
			user.DeletedAt = deletedAtSQL.Time.In(myLocation())
		} else {
			user.DeletedAt = time.Time{} // مقدار زمان صفر برای NULL
		}

		users = append(users, user) // اضافه کردن به slice 'users'
	}
	if err = selectResult.Err(); err != nil {
		return []userDataModel.User{}, 0, fmt.Errorf("خطا در پیمایش نتایج کوئری: %w", err)
	}
	return users, total, nil
}

func (ds *UserDBDS) RenameStudent(ctx context.Context, req userSchema.UpdateUserRequest) (userDataModel.User, error) {
	var students userDataModel.User
	stmt := fmt.Sprintf("UPDATE %s SET  name = ?, family = ?, updated_at = ? WHERE id = ? ", ds.tableName)
	var updatedAt time.Time
	sss, err := ds.db.PrepareContext(ctx, stmt)
	if err != nil {
		return userDataModel.User{}, err
	}
	defer sss.Close()

	result, err := sss.ExecContext(ctx,
		students.Name,
		students.Family,
		updatedAt,
		req.ID,
	)
	if err != nil {
		return userDataModel.User{}, err
	}
	// (optional) require for number of updated column
	rows, err := result.RowsAffected()
	if err != nil {
		return userDataModel.User{}, errors.New("error in number update")
	}
	if rows == 0 {
		return userDataModel.User{}, fmt.Errorf("rows == 0")
	}
	updatedAt = updatedAt.In(myLocation())
	return students, nil

}

func (ds *UserDBDS) UpdateStudent(ctx context.Context, req userSchema.UpdateUserRequest) (userDataModel.User, error) {
	now := time.Now().In(myLocation())
	var students userDataModel.User
	stmt := fmt.Sprintf("UPDATE %s SET updated_at = ? WHERE id = ? ", ds.tableName)
	sss, err := ds.db.PrepareContext(ctx, stmt)
	if err != nil {
		return userDataModel.User{}, err
	}
	defer sss.Close()

	result, err := sss.ExecContext(ctx,
		now,
		req.ID,
	)
	if err != nil {
		return userDataModel.User{}, err
	}
	// (optional) require for number of updated column
	rows, err := result.RowsAffected()
	if err != nil {
		return userDataModel.User{}, errors.New("error in number update")
	}
	if rows == 0 {
		return userDataModel.User{}, fmt.Errorf("rows == 0")
	}

	readQuery := fmt.Sprintf("SELECT id , code , name , family , created_at , updated_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)
	var createdAt, updatedAt, deletedAt sql.NullTime

	if err = ds.db.QueryRowContext(ctx, readQuery, req.ID).Scan(&students.ID, &students.Code, &students.Name, &students.Family, &createdAt, &updatedAt, &deletedAt); err != nil {
		return userDataModel.User{}, err
	}

	if createdAt.Valid {
		students.CreatedAt = createdAt.Time.In(myLocation())
	} else {
		students.CreatedAt = time.Time{}
	}

	if updatedAt.Valid {
		fmt.Println(updatedAt.Time)
		students.UpdatedAt = updatedAt.Time.In(myLocation())
	} else {
		students.UpdatedAt = time.Time{}
	}
	if deletedAt.Valid {
		students.DeletedAt = deletedAt.Time.In(myLocation())
	} else {
		students.DeletedAt = time.Time{}
	}
	return students, nil

}

func (ds *UserDBDS) readTaskByID(ctx context.Context, userID int64) (userDataModel.User, error) {
	var students userDataModel.User

	readQuery := fmt.Sprintf("SELECT id , code , name , family , created_at , updated_at , deleted_at FROM %s WHERE id = ?", ds.tableSQL)
	var createdAt, updatedAt, deletedAt sql.NullTime

	if err := ds.db.QueryRowContext(ctx, readQuery, userID).Scan(&students.ID, &students.Code, &students.Name, &students.Family, &createdAt, &updatedAt, &deletedAt); err != nil {
		return userDataModel.User{}, err
	}

	if createdAt.Valid {
		students.CreatedAt = createdAt.Time.In(myLocation())
	} else {
		students.CreatedAt = time.Time{}
	}

	if updatedAt.Valid {
		fmt.Println(updatedAt.Time)
		students.UpdatedAt = updatedAt.Time.In(myLocation())
	} else {
		students.UpdatedAt = time.Time{}
	}
	if deletedAt.Valid {
		students.DeletedAt = deletedAt.Time.In(myLocation())
	} else {
		students.DeletedAt = time.Time{}
	}

	return students, nil

}
