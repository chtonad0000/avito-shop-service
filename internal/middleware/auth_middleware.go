package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/avito-shop-service/pkg/auth"
)

type key string

const EmployeeUsernameKey key = "username"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := auth.ParseToken(tokenStr)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), EmployeeUsernameKey, claims.EmployeeUsername)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetEmployeeUsername(ctx context.Context) (string, bool) {
	employeeUsername, ok := ctx.Value(EmployeeUsernameKey).(string)
	return employeeUsername, ok
}
