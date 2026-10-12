package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey contextKey = "user"

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func AuthMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, `{"success":false,"message":"Header Authorization tidak ditemukan"}`, http.StatusUnauthorized)
				return
			}

			tokenParts := strings.Split(authHeader, " ")
			if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
				http.Error(w, `{"success":false,"message":"Format token tidak valid (Gunakan Bearer <token>)"}`, http.StatusUnauthorized)
				return
			}

			tokenString := tokenParts[1]
			claims := &Claims{}

			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, `{"success":false,"message":"Token kedaluwarsa atau tidak valid"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(UserContextKey).(*Claims)
			if !ok {
				http.Error(w, `{"success":false,"message":"Sesi pengguna tidak ditemukan"}`, http.StatusUnauthorized)
				return
			}

			allowed := false
			for _, role := range roles {
				if claims.Role == role || claims.Role == "SUPER_ADMIN" {
					allowed = true
					break
				}
			}

			if !allowed {
				http.Error(w, `{"success":false,"message":"Akses ditolak: Anda tidak memiliki wewenang untuk aksi ini"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
