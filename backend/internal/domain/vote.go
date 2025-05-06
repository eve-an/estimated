package domain

import (
	"errors"
	"fmt"
	"time"

	"github.com/eve-an/estimated/internal/value"
)

type VoteEntry struct {
	ID        string
	Session   string
	Name      string
	Value     value.VoteValue
	Timestamp time.Time
}

func NewVoteEntry(id, session, name string, val value.VoteValue, ts time.Time) (VoteEntry, error) {
	v := VoteEntry{
		ID:        id,
		Session:   session,
		Name:      name,
		Value:     val,
		Timestamp: ts,
	}

	if err := v.Valid(); err != nil {
		return VoteEntry{}, err
	}

	return v, nil
}

func (v VoteEntry) Valid() error {
	if !v.Value.IsValid() {
		return fmt.Errorf("invalid vote value: %+v", v.Value)
	}

	if v.Name == "" {
		return errors.New("vote name must not be empty")
	}

	if v.ID == "" {
		return errors.New("vote ID must not be empty")
	}

	if v.Session == "" {
		return errors.New("vote session must not be empty")
	}

	if v.Timestamp.After(time.Now()) {
		return fmt.Errorf("vote timestamp must not be in the future: %+v", v.Timestamp)
	}

	return nil
}
