package httpx

import (
	"context"
	"errors"
)

type contextKey string

const sessionKey contextKey = "session"

var ErrSessionNotFound = errors.New("no session was generated and saved in context")

func SessionKeyFromContext(ctx context.Context) (string, error) {
	key, found := ctx.Value(sessionKey).(string)
	if !found {
		return "", ErrSessionNotFound
	}

	return key, nil
}
