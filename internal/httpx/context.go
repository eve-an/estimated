//go:generate go-enum --marshal
package httpx

import (
	"context"
	"errors"
)

// ENUM(session)
type ContextKey string

var ErrSessionNotFound = errors.New("no session was generated and saved in context")

func SessionKeyFromContext(ctx context.Context) (string, error) {
	key, found := ctx.Value(ContextKeySession).(string)
	if !found {
		return "", ErrSessionNotFound
	}

	return key, nil
}
