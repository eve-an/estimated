package store

import (
	"sync"
	"time"

	"github.com/eve-an/estimated/internal/domain"
	"github.com/eve-an/estimated/internal/infra/collection"
	"github.com/eve-an/estimated/internal/infra/notify"
	"github.com/eve-an/estimated/internal/service"
)

const maxVotingCount = 100

type sessionData struct {
	Token     string
	CreatedAt time.Time
	votes     *collection.RingBuffer[domain.VoteEntry]
}

func (s *sessionData) Push(vs []domain.VoteEntry) {
	if len(vs) > maxVotingCount {
		vs = vs[len(vs)-maxVotingCount:] // dont save everything because it'll kill our memory
	}

	for _, v := range vs {
		s.votes.Push(v)
	}
}

func (s *sessionData) Entries() []domain.VoteEntry {
	return s.votes.GetAll()
}

func newSessionData(token string) *sessionData {
	return &sessionData{
		Token:     token,
		CreatedAt: time.Now(),
		votes:     collection.NewRingBuffer[domain.VoteEntry](maxVotingCount),
	}
}

var _ service.VoteStore = (*SessionStore)(nil)

type SessionStore struct {
	sessions map[string]*sessionData

	mu              sync.RWMutex
	sessionNotifier *notify.SessionNotifier
}

func NewSessionStore(sessionNotifier *notify.SessionNotifier) *SessionStore {
	return &SessionStore{
		sessions:        make(map[string]*sessionData),
		mu:              sync.RWMutex{},
		sessionNotifier: sessionNotifier,
	}
}

func (s *SessionStore) Add(token string, vote domain.VoteEntry) error {
	s.mu.Lock()
	value, found := s.sessions[token]
	if !found {
		value = newSessionData(token)
		s.sessions[token] = value
	}
	s.mu.Unlock()

	value.votes.Push(vote)

	s.sessionNotifier.Notify()

	return nil
}

func (s *SessionStore) Clear() (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	totalDeleted := 0
	for _, sess := range s.sessions {
		totalDeleted += sess.votes.Clear()
	}

	return totalDeleted, nil
}

func (s *SessionStore) Exists(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, ok := s.sessions[token]
	return ok
}

func (s *SessionStore) Get(token string) ([]domain.VoteEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sessions[token].Entries(), nil
}

func (s *SessionStore) List() ([]domain.VoteEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	votings := make([]domain.VoteEntry, 0, len(s.sessions)*maxVotingCount)
	for _, data := range s.sessions {
		votings = append(votings, data.Entries()...)
	}

	return votings, nil
}
