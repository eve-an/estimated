package notify

import (
	"errors"
	"sync"

	"github.com/eve-an/estimated/internal/config"
)

var ErrChannelNotFound = errors.New("channel was not found")

type Notifier interface {
	Notify()
	Subscribe(key string) <-chan struct{}
	Unsubscribe(key string) error
}

type SessionNotifier struct {
	channels    map[string]chan struct{}
	mu          sync.Mutex
	channelSize int
}

func NewSessionNotifier(config *config.Config) *SessionNotifier {
	return &SessionNotifier{
		channels:    make(map[string]chan struct{}),
		channelSize: config.SessionChannelSize,
	}
}

func (s *SessionNotifier) Notify() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, c := range s.channels {
		c <- struct{}{} // may block
	}
}

func (s *SessionNotifier) Subscribe(key string) <-chan struct{} {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := make(chan struct{}, s.channelSize)
	s.channels[key] = c

	return c
}

func (s *SessionNotifier) Unsubscribe(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	c, found := s.channels[key]
	if !found {
		return ErrChannelNotFound
	}

	close(c)
	delete(s.channels, key)

	return nil
}
