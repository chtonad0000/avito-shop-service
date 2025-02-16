//go:build unit
// +build unit

package auth

import (
	"github.com/avito-shop-service/pkg/auth"
	"github.com/golang-jwt/jwt/v4"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	token, err := auth.GenerateToken("testuser")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestParseToken_ValidToken(t *testing.T) {
	token, err := auth.GenerateToken("testuser")
	assert.NoError(t, err)
	claims, err := auth.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, "testuser", claims.EmployeeUsername)
}

func TestParseToken_InvalidToken(t *testing.T) {
	claims, err := auth.ParseToken("invalid.token.value")
	assert.Error(t, err)
	assert.Nil(t, claims)
}

func TestParseToken_ExpiredToken(t *testing.T) {
	claims := auth.Claims{
		EmployeeUsername: "testuser",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(auth.JwtSecretKey)
	assert.NoError(t, err)
	parsedClaims, err := auth.ParseToken(tokenStr)
	assert.Error(t, err)
	assert.Nil(t, parsedClaims)
}
