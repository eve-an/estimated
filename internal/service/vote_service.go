package service

import (
	"context"
	"errors"
	"log/slog"

	"github.com/eve-an/estimated/internal/domain"
	"github.com/eve-an/estimated/internal/infra/notify"
)

type VoteStore interface {
	Add(key string, vote domain.VoteEntry) error
	Get(key string) ([]domain.VoteEntry, error)
	List() ([]domain.VoteEntry, error)
	Clear() (int, error)
}

type VoteService interface {
	AddVotes(ctx context.Context, sessionKey string, votes []domain.VoteEntry) error
	GetAllVotes(ctx context.Context) ([]domain.VoteEntry, error)
	ClearAllVotes(ctx context.Context) (int, error)
}

type voteService struct {
	store    VoteStore
	notifier notify.Notifier
	logger   *slog.Logger
}

func NewVoteService(
	store VoteStore,
	notifier notify.Notifier,
	logger *slog.Logger,
) VoteService {
	return &voteService{
		store:    store,
		notifier: notifier,
		logger:   logger,
	}
}

func (s *voteService) AddVotes(
	ctx context.Context,
	sessionKey string,
	votes []domain.VoteEntry,
) error {
	if sessionKey == "" {
		return errors.New("session key is required")
	}

	if len(votes) == 0 {
		return errors.New("no votes provided")
	}

	for i, vote := range votes {
		if err := vote.Valid(); err != nil {
			return errors.New("invalid vote at index " + string(i) + ": " + err.Error())
		}
	}

	for _, vote := range votes {
		if err := s.store.Add(sessionKey, vote); err != nil {
			s.logger.Error("failed to add vote",
				"session_key", sessionKey,
				"vote_timestamp", vote.Timestamp,
				"vote_value", vote.Value,
				"error", err)
			return err
		}
	}

	s.notifier.Notify()

	return nil
}

func (s *voteService) GetAllVotes(_ context.Context) ([]domain.VoteEntry, error) {
	votes, err := s.store.List()
	if err != nil {
		s.logger.Error("failed to retrieve votes", "error", err)
		return nil, err
	}
	return votes, nil
}

func (s *voteService) ClearAllVotes(ctx context.Context) (int, error) {
	count, err := s.store.Clear()
	if err != nil {
		s.logger.Error("failed to clear votes", "error", err)
		return 0, err
	}

	s.notifier.Notify()

	return count, nil
}
