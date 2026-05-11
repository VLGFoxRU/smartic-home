package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserIDKey contextKey = "userID"
const UserRoleKey contextKey = "role"

// AuthMiddleware проверяет JWT в заголовке Authorization и добавляет userID/role в контекст
func AuthMiddleware(secret []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, `{"error":"invalid authorization format"}`, http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]

			claims := &jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return secret, nil
			})
			if err != nil || !token.Valid {
				http.Error(w, `{"error":"invalid or expired token"}`, http.StatusUnauthorized)
				return
			}

			// Добавляем userID и роль в контекст
			ctx := context.WithValue(r.Context(), UserIDKey, (*claims)["user_id"])
			ctx = context.WithValue(ctx, UserRoleKey, (*claims)["role"])
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ValidateToken проверяет JWT и возвращает claims (userID, role) или ошибку
func ValidateToken(tokenString string, secret []byte) (string, string, error) {
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !token.Valid {
		return "", "", err
	}
	userID, _ := (*claims)["user_id"].(string)
	role, _ := (*claims)["role"].(string)
	return userID, role, nil
}