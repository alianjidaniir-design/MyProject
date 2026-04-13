package mysqlDS

import (
	"MyProject/apiSchema/teacherSchema"
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
