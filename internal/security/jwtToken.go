// Package security provides functionalities to work with JWT tokens.
package security

import (
	"errors"
	"time"

	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

// JWTClaims defines the structure of JWT claims used in the token.
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

// Errors related to token processing
var (
	ErrorParseClaims  = errors.New("couldn't parse claims")
	ErrorTokenExpired = errors.New("token expired")
)

// GenerateToken creates a JWT token with the specified user ID.
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

// ValidateToken verifies the validity of a JWT token and checks if it's expired.
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

// GetUserIDFromToken extracts the user ID from a JWT token.
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
