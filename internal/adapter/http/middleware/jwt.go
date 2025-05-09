package middleware

import (
	"context"
	"net/http"
	"strings"

	"ip_detector/internal/auth"
	"ip_detector/internal/logger"
)

type contextKey string

const userEmailKey contextKey = "user_email"

func JWTMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logger.Log.Sugar()

			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				log.Warn("authorization header missing")
				http.Error(w, "Authorization header missing", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				log.Warnw("authorization header without Bearer prefix", "header", authHeader)
				http.Error(w, "Invalid token format", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
			email, err := auth.ParseToken(tokenStr, secret)
			if err != nil {
				log.Warnw("invalid token", "error", err)
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			log.Infow("token verified", "email", email)
			ctx := context.WithValue(r.Context(), userEmailKey, email)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
