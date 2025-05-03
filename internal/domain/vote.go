package domain

import (
	"fmt"
	"time"

	"github.com/eve-an/estimated/internal/value"
)

type VoteEntry struct {
	ID        string
	Session   string
	Value     value.VoteValue
	Timestamp time.Time
}

func (v VoteEntry) Valid() error {
	if !v.Value.IsValid() {
		return fmt.Errorf("invalid vote value: %+v", v.Value)
	}

	return nil
}
