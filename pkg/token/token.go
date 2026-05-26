package token

import (
	tokenDataModel "MyProject/models/token/dataModel"
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

	claim := tokenDataModel.AccessToken{
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

func GenerateRefreshToken(roleID int64, UserID int64) (string, error) {
	exp := time.Now().In(myLocation()).Add(constants.RefreshTokenExpiry)
	var jwtSecret = []byte(configs.RefreshTokenSecret)

	claim := jwt.MapClaims{
		"role_id": roleID,
		"user_id": UserID,
		"scope":   "refresh",
		"exp":     jwt.NewNumericDate(exp),
		"iat":     jwt.NewNumericDate(time.Now().In(myLocation())),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtSecret)

}
