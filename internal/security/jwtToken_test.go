package security

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGenerateToken(t *testing.T) {
	userID := uuid.New()
	token, err := GenerateToken(userID)
	assert.NoError(t, err, "Token generation should not produce an error")
	assert.NotEmpty(t, token, "Token should not be empty")
}

func TestValidateToken(t *testing.T) {
	userID := uuid.New()
	token, _ := GenerateToken(userID)

	err := ValidateToken(token)
	assert.NoError(t, err, "Valid token should not produce an error")

}

func TestGetUserIDFromToken(t *testing.T) {
	userID := uuid.New()
	token, _ := GenerateToken(userID)

	extractedUserID, err := GetUserIDFromToken(token)
	assert.NoError(t, err, "Should not produce an error for a valid token")
	assert.Equal(t, userID, extractedUserID, "Extracted user ID should match the original")
}
