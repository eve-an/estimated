// Code generated by go-enum DO NOT EDIT.
// Version: 0.6.1
// Revision: a6f63bddde05aca4221df9c8e9e6d7d9674b1cb4
// Build Date: 2025-03-18T23:42:14Z
// Built By: goreleaser

package session

import (
	"errors"
	"fmt"
)

const (
	// ContextKeySession is a ContextKey of type session.
	ContextKeySession ContextKey = "session"
)

var ErrInvalidContextKey = errors.New("not a valid ContextKey")

// String implements the Stringer interface.
func (x ContextKey) String() string {
	return string(x)
}

// IsValid provides a quick way to determine if the typed value is
// part of the allowed enumerated values
func (x ContextKey) IsValid() bool {
	_, err := ParseContextKey(string(x))
	return err == nil
}

var _ContextKeyValue = map[string]ContextKey{
	"session": ContextKeySession,
}

// ParseContextKey attempts to convert a string to a ContextKey.
func ParseContextKey(name string) (ContextKey, error) {
	if x, ok := _ContextKeyValue[name]; ok {
		return x, nil
	}
	return ContextKey(""), fmt.Errorf("%s is %w", name, ErrInvalidContextKey)
}
