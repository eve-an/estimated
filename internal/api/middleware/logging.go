package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/eve-an/estimated/internal/infra/session"
)

type Middleware struct {
	logger *slog.Logger
}

func NewMiddleware(logger *slog.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

func (m *Middleware) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		duration := time.Since(start)

		m.logger.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"duration", duration,
			"remote", r.RemoteAddr,
			"user_agent", r.UserAgent(),
		)
	})
}

func (m *Middleware) AddSessionCookie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookieKey := session.ContextKeySession.String()

		cookie, err := r.Cookie(cookieKey)
		if err != nil {
			cookie = &http.Cookie{
				Name:     cookieKey,
				Value:    generateToken(),
				Path:     "/",
				HttpOnly: true,
				Secure:   false, // Set to true in production (with HTTPS)
				SameSite: http.SameSiteLaxMode,
				MaxAge:   86400 * 3, // 3 days
			}
		}

		http.SetCookie(w, cookie)

		ctx := context.WithValue(r.Context(), session.ContextKeySession, cookie.Value)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
