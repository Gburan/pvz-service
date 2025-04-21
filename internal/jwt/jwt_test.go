package jwt

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateAndParseToken(t *testing.T) {
	secret := "supersecretkey"
	role := "employee"
	uid := "user-123"
	expiresIn := time.Hour

	t.Run("Generate and parse token successfully", func(t *testing.T) {
		tokenStr, err := GenerateToken(secret, role, uid, expiresIn)
		require.NoError(t, err)
		require.NotEmpty(t, tokenStr)

		parsedRole, err := ParseToken(tokenStr, secret)
		require.NoError(t, err)
		assert.Equal(t, role, parsedRole)
	})

	t.Run("ParseToken with wrong secret", func(t *testing.T) {
		validToken, err := GenerateToken(secret, role, uid, expiresIn)
		require.NoError(t, err)

		parsedRole, err := ParseToken(validToken, "wrongsecret")
		assert.Error(t, err)
		assert.Empty(t, parsedRole)
	})

	t.Run("ParseToken with invalid token format", func(t *testing.T) {
		parsedRole, err := ParseToken("this.is.not.a.token", secret)
		assert.Error(t, err)
		assert.Empty(t, parsedRole)
	})

	t.Run("ParseToken with unexpected signing method", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
			"role": role,
			"id":   uid,
			"iat":  jwt.NewNumericDate(time.Now()),
			"exp":  jwt.NewNumericDate(time.Now().Add(expiresIn)),
		})
		tokenString, _ := token.SignedString([]byte("dummy"))

		parsedRole, err := ParseToken(tokenString, secret)
		assert.Error(t, err)
		assert.Empty(t, parsedRole)
	})

	t.Run("ParseToken with missing role claim", func(t *testing.T) {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"id":  uid,
			"iat": jwt.NewNumericDate(time.Now()),
			"exp": jwt.NewNumericDate(time.Now().Add(expiresIn)),
		})

		tokenString, err := token.SignedString([]byte(secret))
		require.NoError(t, err)

		parsedRole, err := ParseToken(tokenString, secret)
		assert.Error(t, err)
		assert.Empty(t, parsedRole)
	})
}
