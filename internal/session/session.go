package session

import (
	"sync"
	"time"

	"github.com/eve-an/estimated/internal/collection"
	"github.com/eve-an/estimated/internal/db"
	"github.com/eve-an/estimated/internal/model"
	"github.com/eve-an/estimated/internal/notify"
)

const maxVotingCount = 100

type sessionData struct {
	Token     string
	CreatedAt time.Time
	votes     *collection.RingBuffer[model.VoteEntry]
}

func (s *sessionData) Push(vs []model.VoteEntry) {
	if len(vs) > maxVotingCount {
		vs = vs[len(vs)-maxVotingCount:] // dont save everything because it'll kill our memory
	}

	for _, v := range vs {
		s.votes.Push(v)
	}
}

func (s *sessionData) Entries() []model.VoteEntry {
	return s.votes.GetAll()
}

func newSessionData(token string) *sessionData {
	return &sessionData{
		Token:     token,
		CreatedAt: time.Now(),
		votes:     collection.NewRingBuffer[model.VoteEntry](maxVotingCount),
	}
}

var _ db.VoteEntryStore = (*SessionStore)(nil)

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

func (s *SessionStore) Add(token string, vote model.VoteEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, found := s.sessions[token]
	if !found {
		value = newSessionData(token)
		s.sessions[token] = value
	}

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

func (s *SessionStore) Get(token string) ([]model.VoteEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.sessions[token].Entries(), nil
}

func (s *SessionStore) List() ([]model.VoteEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	votings := make([]model.VoteEntry, 0, len(s.sessions)*maxVotingCount)
	for _, data := range s.sessions {
		votings = append(votings, data.Entries()...)
	}

	return votings, nil
}
