package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/eve-an/estimated/internal/api"
	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/api/mapper"
	"github.com/eve-an/estimated/internal/domain"
	"github.com/eve-an/estimated/internal/infra/session"
	"github.com/eve-an/estimated/internal/service"
	"github.com/go-chi/chi/v5"
)

type VotesHandler struct {
	logger *slog.Logger

	voteService   service.VoteService
	voteMapper    *mapper.VoteMapper
	nameGenerator service.NameGenerator
}

func NewVotesHandler(
	logger *slog.Logger,
	voteService service.VoteService,
	voteMapper *mapper.VoteMapper,
	nameGenerator service.NameGenerator,
) *VotesHandler {
	return &VotesHandler{
		logger:        logger,
		voteService:   voteService,
		voteMapper:    voteMapper,
		nameGenerator: nameGenerator,
	}
}

func (s *VotesHandler) Add(w http.ResponseWriter, r *http.Request) {
	var requestDTO dto.VoteRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&requestDTO); err != nil {
		s.logger.Error("reading body failed", "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "failed to read body",
		})
		return
	}

	key, err := session.FromContext(r.Context())
	if err != nil {
		s.logger.Warn("client has no session key", "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "no session found",
		})
		return
	}

	name, err := s.nameGenerator.NameFor(key)
	if err != nil {
		s.logger.Warn("failed to generate name", "err", err, "session_key", key)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "failed to generate name",
		})
		return
	}

	domainVote, err := s.voteMapper.RequestToDomain(&requestDTO, key, name)
	if err != nil {
		api.WriteJSON(w, http.StatusBadRequest, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "invalid vote",
			Data:   err.Error(),
		})
		return
	}

	slog.Debug("vote request ",
		"request", requestDTO,
		"model", domainVote,
	)

	if err := s.voteService.AddVotes(r.Context(), key, []domain.VoteEntry{domainVote}); err != nil {
		s.logger.Error("failed to add vote", "err", err, "vote", domainVote)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "failed to store vote",
		})
		return
	}

	api.WriteJSON(w, http.StatusOK, dto.APIResponse{
		Status: dto.StatusSuccess,
		Data:   domainVote,
	})
}

func (s *VotesHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	votes, err := s.voteService.GetAllVotes(r.Context())
	if err != nil {
		s.logger.Error("could not fetch all votes from the store", "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "could not fetch all votes from the store",
		})
		return
	}

	response := s.voteMapper.DomainToResponse(votes)
	s.logger.DebugContext(r.Context(), "all vote entries",
		"entries", votes,
		"response", response,
	)

	api.WriteJSON(w, http.StatusOK, dto.APIResponse{
		Status: dto.StatusSuccess,
		Data:   response,
	})
}

func (s *VotesHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	n, err := s.voteService.ClearAllVotes(r.Context())
	if err != nil {
		s.logger.Error("could not delete all votes from the store", "err", err)
		api.WriteJSON(w, http.StatusInternalServerError, dto.APIResponse{
			Status: dto.StatusError,
			Error:  "could not delete all votes from store",
		})
		return
	}

	api.WriteJSON(w, http.StatusOK, dto.APIResponse{
		Status: dto.StatusSuccess,
		Data:   n,
	})
}

func (s *VotesHandler) Routes() http.Handler {
	r := chi.NewRouter()

	r.Post("/", s.Add)
	r.Get("/", s.GetAll)
	r.Delete("/", s.DeleteAll)

	return r
}
