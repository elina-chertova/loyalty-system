package security

import (
	"errors"
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

var (
	ErrorParseClaims  = errors.New("couldn't parse claims")
	ErrorTokenExpired = errors.New("token expired")
)

func GenerateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.TokenExp)),
			},
			UserID: userID,
		},
	)

	tokenString, err := token.SignedString([]byte(config.SecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(signedToken string) error {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		signedToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		},
	)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return ErrorParseClaims
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return ErrorTokenExpired
	}
	return nil
}

func GetUserIDFromToken(signedToken string) (uuid.UUID, error) {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		signedToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(config.SecretKey), nil
		},
	)
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return uuid.Nil, ErrorParseClaims
	}
	return claims.UserID, nil
}
