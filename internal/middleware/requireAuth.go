package middleware

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ab91dev/codeolx/internal/httpx"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	authorization = "Authorization"
)

type AuthClaims struct {
	jwt.RegisteredClaims
}

func RequireAuth(logger *slog.Logger, secret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := logger.With("request_id", RequestIDFromContext(ctx))
			authHeader := r.Header.Get(authorization)
			if authHeader == "" {
				log.Error("no authorization found")
				httpx.Error(w, http.StatusUnauthorized, "auth is required", httpx.CodeUnauthenticated)
				return
			}

			raw, ok := strings.CutPrefix(authHeader, "Bearer ")
			if !ok || raw == "" {
				log.Error("no jwt token found")
				httpx.Error(w, http.StatusUnauthorized, "auth is required", httpx.CodeUnauthenticated)
				return
			}

			var claims AuthClaims
			_, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
				if _, xyz := t.Method.(*jwt.SigningMethodHMAC); !xyz {
					return nil, errors.New("unexpected signing method")
				}

				return secret, nil
			}, jwt.WithValidMethods([]string{"HS256"}))
			if err != nil {
				log.Info("token rejected", "error", err)
				httpx.Error(w, http.StatusUnauthorized, "token expired or invalid", httpx.CodeUnauthenticated)
				return
			}

			userId, err := uuid.Parse(claims.Subject)
			if err != nil {
				log.Info("string to uuid failed", "error", err)
				httpx.Error(w, http.StatusUnauthorized, "token expired or invalid", httpx.CodeUnauthenticated)
				return
			}

			ctxWithUserID := context.WithValue(ctx, userIDKey, userId)
			next.ServeHTTP(w, r.WithContext(ctxWithUserID))
		})
	}
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}
