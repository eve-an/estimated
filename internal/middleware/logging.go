package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/eve-an/estimated/internal/session"
)

type contextKey string

const sessionKey contextKey = "session"

type Middleware struct {
	logger   *slog.Logger
	sessions *session.SessionStore
}

func NewMiddleware(logger *slog.Logger, sessions *session.SessionStore) *Middleware {
	return &Middleware{
		logger:   logger,
		sessions: sessions,
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
		cookie, err := r.Cookie(string(sessionKey))
		if err != nil {
			cookie = &http.Cookie{
				Name:     string(sessionKey),
				Value:    "",
				Path:     "/",
				HttpOnly: true,
				Secure:   false, // Set to true in production (with HTTPS)
				SameSite: http.SameSiteLaxMode,
				MaxAge:   86400 * 3, // 3 days
			}
		}

		if cookie.Value == "" {
			cookie.Value = generateToken()
		}

		http.SetCookie(w, cookie)

		m.sessions.Create(cookie.Value)

		ctx := context.WithValue(r.Context(), sessionKey, cookie.Value)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func generateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
