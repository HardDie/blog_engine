package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/HardDie/blog_engine/internal/service/auth"
	"github.com/HardDie/blog_engine/internal/utils"
)

type AuthMiddleware struct {
	authService auth.IAuth
}

func NewAuthMiddleware(authService auth.IAuth) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}
func (m *AuthMiddleware) RequestMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie := utils.GetCookie(r)

		// If we got no cookie
		if cookie == nil || cookie.Value == "" {
			http.Error(w, "Session not found in cookie", http.StatusUnauthorized)
			return
		}

		// Validate if cookie is active
		ctx := r.Context()
		session, err := m.authService.ValidateCookie(ctx, cookie.Value)
		if err != nil || session == nil {
			switch {
			case errors.Is(err, auth.ErrorSessionNotFound):
				http.Error(w, "Session not found", http.StatusUnauthorized)
				return
			case errors.Is(err, auth.ErrorSessionHasExpired):
				http.Error(w, "Session has expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, "userID", session.UserID)
		ctx = context.WithValue(ctx, "session", session)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
