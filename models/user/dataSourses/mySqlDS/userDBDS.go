package mySqlDS

import (
	"MyProject/apiSchema/userSchema"
	userDataModel "MyProject/models/user/dataModel"
	"context"
	"fmt"
)

type UserDBDS struct {
	tablename string
	tableSQL  string
	db        DBExecture
}

func NewUserDBDSFromEnv() (*UserDBDS, bool, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return nil, false, err
	}
	tableSQL, err := studentTableName(cfg.StudentTableName)
	if err != nil {
		return nil, false, err
	}
	db, err := Open(cfg)
	if err != nil {
		return nil, false, err
	}
	if err := EnsureStudentTable(db, cfg.StudentTableName); err != nil {
		return nil, false, err
	}
	return &UserDBDS{
		tablename: cfg.StudentTableName,
		tableSQL:  tableSQL,
		db:        db,
	}, true, nil
}
func (ds *UserDBDS) CreateStudent(ctx context.Context, req userSchema.LoginRequest) (userDataModel.User, error) {
	insertQuery := fmt.Sprintf("INSERT INTO %s (code , name , family ) VALUES (?, ? , ?)", ds.tableSQL)
	insertResult, err := ds.db.ExecContext(ctx, insertQuery, req.Code, req.Name, req.Family)
	if err != nil {
		return userDataModel.User{}, err
	}
	insertedID, err := insertResult.LastInsertId()
	if err != nil {
		return userDataModel.User{}, err
	}
	return ds.readTaskByID(ctx, insertedID)
}

func (ds *UserDBDS) readTaskByID(ctx context.Context, userID int64) (userDataModel.User, error) {
	var students userDataModel.User
	readQuery := fmt.Sprintf("SELECT id , code , name , family FROM %s WHERE id = ?", ds.tableSQL)
	if err := ds.db.QueryRowContext(ctx, readQuery, userID).Scan(&students.ID, &students.Code, students.Name, students.Family); err != nil {
		return userDataModel.User{}, err
	}
	return students, nil

}
