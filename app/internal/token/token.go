package token

import (
	"auth-service/app/internal/config"
	"auth-service/app/internal/exception"
	"auth-service/app/internal/lib"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type Claims struct {
	UserID   uint     `json:"user_id"`
	UserRole lib.Role `json:"role"`
	jwt.RegisteredClaims
}

func CreateRefresh() string {
	return uuid.New().String()
}

func CreateAccess(userID uint, userRole lib.Role) (string, error) {
	config := config.Get()

	issuedAt := time.Now().UTC()
	expiresTime := issuedAt.Add(1 * time.Hour)
	claims := Claims{
		UserID:   userID,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenSigned, err := token.SignedString([]byte(config.Token.SecretKey))
	if err != nil {
		return "", exception.ErrInternal
	}

	return tokenSigned, nil
}

func VerifyAccess(tokenSigned string) (*Claims, error) {
	config := config.Get()

	token, err := jwt.ParseWithClaims(
		tokenSigned,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, exception.ErrVerify
			}
			return config.Token.SecretKey, nil
		},
	)
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, exception.ErrVerify
	}

	return claims, nil

}
