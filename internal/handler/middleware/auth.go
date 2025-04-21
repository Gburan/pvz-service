package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"pvz-service/internal/handler"
	"pvz-service/internal/jwt"
)

var (
	ErrNoToken          = errors.New("no token provided")
	ErrNoAcceptableRole = errors.New("no acceptable role")
)

type UserRole string

const (
	Employee  UserRole = "EMPLOYEE"
	Moderator UserRole = "MODERATOR"

	authorisationPrefix = "Bearer "
	roleKey             = "role"
)

func hasRequiredRole(userRole UserRole, requiredRoles []UserRole) bool {
	role := UserRole(strings.ToUpper(string(userRole)))
	for _, requiredRole := range requiredRoles {
		if role == requiredRole {
			return true
		}
	}
	return false
}

func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoToken
	}

	return strings.TrimPrefix(authHeader, authorisationPrefix), nil
}

func AuthMiddleware(secret string, requiredRoles []UserRole, next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := extractToken(r)
		if err != nil {
			handler.RespondWithError(w, http.StatusUnauthorized, "authorization required", err)
			return
		}

		role, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			handler.RespondWithError(w, http.StatusUnauthorized, "invalid token", err)
			return
		}

		userRole := UserRole(role)
		if !hasRequiredRole(userRole, requiredRoles) {
			handler.RespondWithError(w, http.StatusForbidden, "insufficient permissions", ErrNoAcceptableRole)
			return
		}

		ctx := context.WithValue(r.Context(), roleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
