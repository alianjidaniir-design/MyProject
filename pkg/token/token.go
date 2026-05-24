package token

import (
	studentDataModel "MyProject/models/student/dataModel"
	"MyProject/statics/configs"
	"MyProject/statics/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func myLocation() *time.Location {
	loc, err := time.LoadLocation("Asia/ُTehran")
	if err != nil {
		return time.FixedZone("Asia/Tehran", 3*3600+30*60)
	}
	return loc
}

func GenerateAccessToken(userID int64, roleID int64) (string, error) {
	var jwtSecret = []byte(configs.AccessTokenSecret)
	jti := uuid.NewString()

	exp := time.Now().In(myLocation()).Add(constants.AccessTokenExpiry)

	claim := studentDataModel.AccessToken{
		UserID: userID,
		RoleID: roleID,
		Scope:  "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(myLocation())),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtSecret)

}

func GenerateRefreshToken() (string, error) {
	exp := time.Now().In(myLocation()).Add(constants.AccessTokenExpiry)
	var jwtSecret = []byte(configs.AccessTokenSecret)

	claim := studentDataModel.AccessToken{
		Scope: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(myLocation())),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtSecret)

}
