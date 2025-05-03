package service

import (
	"context"
	"log/slog"

	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/infra/notify"
)

type EventService interface {
	Subscribe(ctx context.Context, sessionKey string) <-chan *dto.EventResponseDTO
}

type eventService struct {
	logger          *slog.Logger
	sessionNotifier *notify.SessionNotifier
	voteService     VoteService
}

func NewEventService(
	logger *slog.Logger,
	sessionNotifier *notify.SessionNotifier,
	voteService VoteService,
) EventService {
	return &eventService{
		logger:          logger,
		sessionNotifier: sessionNotifier,
		voteService:     voteService,
	}
}

func (e *eventService) Subscribe(ctx context.Context, sessionKey string) <-chan *dto.EventResponseDTO {
	responseChannel := make(chan *dto.EventResponseDTO, 10)
	notification := e.sessionNotifier.Subscribe(sessionKey)

	go func() {
		for {
			select {
			case <-ctx.Done():

				responseChannel <- nil

			case <-notification:
				votes, err := e.voteService.GetAllVotes(ctx)
				if err != nil {
					e.logger.Error("could not extract votes from store", "err", err)
					continue
				}

				responseDTO := dto.EventResponseDTO{
					Name:  sessionKey, // todo map name
					Votes: votes,
				}

				responseChannel <- &responseDTO
			}
		}
	}()

	return responseChannel
}
