package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrFailedSign      = errors.New("failed to sign token")
	ErrFailedParse     = errors.New("failed to parse token")
	ErrUnexpSignMethod = errors.New("unexpected signing method")
)

func GenerateToken(secret, role, uid string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"role": role,
		"id":   uid,
		"iat":  jwt.NewNumericDate(time.Now()),
		"exp":  jwt.NewNumericDate(time.Now().Add(expiresIn)),
	})

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrFailedSign, err)
	}
	return tokenString, nil
}

func ParseToken(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrUnexpSignMethod, token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrFailedParse, err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, ok := claims["role"].(string)
		if !ok {
			return "", fmt.Errorf("invalid token: role claim missing or not a string")
		}
		return role, nil
	}

	return "", jwt.ErrInvalidKey
}
