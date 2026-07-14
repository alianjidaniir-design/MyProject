package token

import (
	tokenDataModel "MyProject/models/token/dataModel"
	"MyProject/pkg/timeLoc"
	"MyProject/statics/configs"
	"MyProject/statics/constants"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func GenerateAccessToken(userID int64, roleName string) (string, error) {
	var jwtSecret = []byte(configs.AccessTokenSecret)
	jti := uuid.NewString()

	exp := time.Now().In(timeLoc.MyLocation()).Add(constants.AccessTokenExpiry)

	claim := tokenDataModel.AccessToken{
		UserID:   userID,
		RoleName: roleName,
		Scope:    "access",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now().In(timeLoc.MyLocation())),
			ID:        jti,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtSecret)

}

func GenerateRefreshToken(roleName string, UserID int64) (string, error) {
	exp := time.Now().In(timeLoc.MyLocation()).Add(constants.RefreshTokenExpiry)
	var jwtSecret = []byte(configs.RefreshTokenSecret)

	claim := jwt.MapClaims{
		"role_name": roleName,
		"user_id":   UserID,
		"scope":     "refresh",
		"exp":       jwt.NewNumericDate(exp),
		"iat":       jwt.NewNumericDate(time.Now().In(timeLoc.MyLocation())),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString(jwtSecret)

}
