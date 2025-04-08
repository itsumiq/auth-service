package token

import (
	"auth-service/app/internal/exception"
	role "auth-service/app/internal/lib"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var SECRET_KEY = []byte("qeqewqew")

type Claims struct {
	UserID   uint      `json:"user_id"`
	UserRole role.Role `json:"role"`
	jwt.RegisteredClaims
}

func CreateRefresh() string {
	return uuid.New().String()
}

func CreateAccess(userID uint, userRole role.Role) (string, error) {
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
	tokenSigned, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", exception.ErrInternal
	}

	return tokenSigned, nil
}

func VerifyAccess(tokenSigned string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(
		tokenSigned,
		&Claims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, exception.ErrVerify
			}
			return SECRET_KEY, nil
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
