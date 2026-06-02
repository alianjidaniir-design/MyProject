package mysqlDS

import (
	"MyProject/apiSchema/adminSchema"
	adminDataModels "MyProject/models/admins/dataModels"
	tokenDataModel "MyProject/models/token/dataModel"
	"MyProject/pkg/hash"
	TimeLoc "MyProject/pkg/timeLoc"
	"MyProject/pkg/token"
	Val "MyProject/pkg/val"
	"MyProject/statics/constants"
	"fmt"

	"MyProject/models/admins/dataSources"
	"MyProject/models/teachers/dataModels"
	"context"
	"database/sql"
	"errors"
	"time"
)

type AdminDBDS struct {
	tableName string
	db        *sql.DB
}

func NewAdminDBDS(tableName string, db *sql.DB) (dataSources.AdminDS, error) {
	if db == nil {
		return nil, errors.New("db is nil")
	}
	ff := &AdminDBDS{
		tableName: tableName,
		db:        db,
	}

	return ff, nil

}

func (ds *AdminDBDS) CreateAdmin(ctx context.Context, req adminSchema.InformationSchema) (res adminDataModels.Admins, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return adminDataModels.Admins{}, err
	}
	hashing, err := hash.HashingPassword(req.Password)
	if err != nil {
		return adminDataModels.Admins{}, err
	}
	var check bool
	checkRole := `
SELECT
CASE WHEN EXISTS (SELECT 1 FROM roles WHERE name = ?) THEN 1 ELSE 0 END`
	err = ds.db.QueryRowContext(ctx, checkRole, req.RoleName).Scan(&check)
	if err != nil {
		return adminDataModels.Admins{}, err
	}
	if !check {
		return adminDataModels.Admins{}, errors.New("role does not exist")
	}

	now := time.Now().In(TimeLoc.MyLocation())

	insert := fmt.Sprintf("INSERT INTO %s (user_name , password , name , family , email , role_name , created_at ) VALUES (?, ?, ?, ?, ?, ?, ?)", ds.tableName)
	insertQuery, err := ds.db.ExecContext(ctx, insert, req.Username, hashing, req.Name, req.Family, req.Email, req.RoleName, now)
	if err != nil {
		return adminDataModels.Admins{}, err
	}

	insertID, err := insertQuery.LastInsertId()
	if err != nil {
		return adminDataModels.Admins{}, err
	}
	return ds.readQuery(ctx, insertID)

}

func (ds *AdminDBDS) Login(ctx context.Context, req adminSchema.LoginAdminRequest) (access string, refresh string, massage string, err error) {
	err = Val.CheckValidation(req)
	if err != nil {
		return "", "", "", err
	}
	student, err := ds.checkAdmin(ctx, req.Username)
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

func (ds *AdminDBDS) Refresh(ctx context.Context, req string) (access string, refresh string, err error) {
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
	checkToken := fmt.Sprintf("SELECT user_id, role_name, token, expires_at, rekoved_at FROM refreshs WHERE token = ? AND rekoved_at = false")
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

func (ds *AdminDBDS) Logout(ctx context.Context, req adminSchema.LogoutSchema, ref string) error {
	err := Val.CheckValidation(req)
	if err != nil {
		return err
	}
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
	if ref == "" {
		return errors.New("token is empty")
	}
	checkToken := fmt.Sprintf("SELECT user_id, role_name, token, expires_at, rekoved_at FROM refreshs WHERE token = ? AND rekoved_at = false")
	err = tx.QueryRowContext(ctx, checkToken, ref).Scan(&rt.UserID, &rt.RoleName, &rt.Token, &rt.ExpiresAt, &rt.RevokedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) == true {
			return errors.New("Refresh token is invalid")
		}
		return err

	} else if time.Now().In(TimeLoc.MyLocation()).After(rt.ExpiresAt) {
		return errors.New("token expired")
	}

	updated := fmt.Sprintf("UPDATE refreshs SET rekoved_at = true WHERE token = ?")
	rows, err := tx.PrepareContext(ctx, updated)
	if err != nil {
		return err
	}
	defer rows.Close()
	row, err := rows.ExecContext(ctx, rt.Token)
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

func (ds *AdminDBDS) checkAdmin(ctx context.Context, userName string) (adminDataModels.Admins, error) {
	var admin adminDataModels.Admins
	checkQuery := "SELECT * FROM admins WHERE user_name = ?"
	err := ds.db.QueryRowContext(ctx, checkQuery, userName).Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Name, &admin.Family, &admin.Email, &admin.RoleName, &admin.CreatedAt)
	if err != nil {
		return admin, err
	}
	return admin, nil
}

func (ds *AdminDBDS) readQuery(ctx context.Context, ID int64) (adminDataModels.Admins, error) {
	var admin adminDataModels.Admins
	read := fmt.Sprintf("SELECT ID , user_name , password , name, family , email , role_name ,  created_at FROM %s WHERE ID=? ", ds.tableName)

	err := ds.db.QueryRowContext(ctx, read, ID).Scan(&admin.ID, &admin.Username, &admin.Password, &admin.Name, &admin.Family, &admin.Email, &admin.RoleName, &admin.CreatedAt)
	if err != nil {
		return admin, err
	}

	return admin, nil
}

func (ds *AdminDBDS) checkingAdmin(ctx context.Context, email string, code string) (dataModels.Teacher, error) {
	var teacher dataModels.Teacher
	checkQuery := "SELECT * FROM teachers WHERE email = ? AND national_code = ?"
	err := ds.db.QueryRowContext(ctx, checkQuery, email, code).Scan(&teacher.ID, &teacher.Name, &teacher.LastName, &teacher.RoleName, &teacher.NationalCode, &teacher.Email, &teacher.Phone, &teacher.WorkExperience, &teacher.Password, &teacher.CreatedAt, &teacher.UpdatedAt, &teacher.DeletedAt)
	if err != nil {
		return teacher, err
	}
	return teacher, nil
}
