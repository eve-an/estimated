package session

import (
	"sync"
	"time"

	"github.com/eve-an/estimated/internal/collection"
)

const maxVotingCount = 100

type Voting struct {
	Value     int       `json:"value"`
	TimeStamp time.Time `json:"timestamp"`
}

type SessionData struct {
	Token     string
	CreatedAt time.Time
	votes     *collection.RingBuffer[Voting]
	updater   chan struct{}
}

func (s *SessionData) Push(vs []Voting) {
	if len(vs) > maxVotingCount {
		vs = vs[len(vs)-maxVotingCount:] // dont save everything because it'll kill our memory
	}

	for _, v := range vs {
		s.votes.Push(v)
	}

	s.updater <- struct{}{}
}

func (s *SessionData) GetVotings() []Voting {
	return s.votes.GetAll()
}

func newSessionData(token string, updater chan struct{}) *SessionData {
	return &SessionData{
		Token:     token,
		CreatedAt: time.Now(),
		votes:     collection.NewRingBuffer[Voting](maxVotingCount),
		updater:   updater,
	}
}

type SessionStore struct {
	sessions map[string]*SessionData

	Updater chan struct{}

	mu sync.RWMutex
}

func NewSessionStore() *SessionStore {
	return &SessionStore{
		sessions: make(map[string]*SessionData),
		mu:       sync.RWMutex{},
		Updater:  make(chan struct{}),
	}
}

func (s *SessionStore) Create(token string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.sessions[token]; !found {
		s.sessions[token] = newSessionData(token, s.Updater)
	}
}

func (s *SessionStore) DeleteAll() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	totalDeleted := 0
	for _, sess := range s.sessions {
		totalDeleted += sess.votes.Clear()
	}

	return totalDeleted
}

func (s *SessionStore) Exists(token string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.sessions[token]
	return ok
}

func (s *SessionStore) Get(token string) *SessionData {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.sessions[token]
}

func (s *SessionStore) GetAllVotings() []Voting {
	s.mu.RLock()
	defer s.mu.RUnlock()

	votings := make([]Voting, 0, len(s.sessions)*maxVotingCount)
	for _, data := range s.sessions {
		votings = append(votings, data.GetVotings()...)
	}

	return votings
}
