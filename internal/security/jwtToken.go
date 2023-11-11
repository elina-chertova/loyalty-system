package security

import (
	"github.com/elina-chertova/loyalty-system/internal/config"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
)

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const TOKEN_EXP = time.Minute * 3
const SECRET_KEY = "supersecretkey"

func GenerateToken(userID uuid.UUID) (string, error) {
	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, JWTClaims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
			},
			UserID: userID,
		},
	)

	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(signedToken string) error {
	claims := &JWTClaims{}
	token, err := jwt.ParseWithClaims(
		signedToken, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		return err
	}
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return config.ErrorParseClaims
	}
	if claims.ExpiresAt.Unix() < time.Now().Local().Unix() {
		return config.ErrorTokenExpired
	}
	return nil
}
