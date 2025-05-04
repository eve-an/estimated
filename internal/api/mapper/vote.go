package mapper

import (
	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/domain"
	"github.com/eve-an/estimated/internal/value"
	"github.com/google/uuid"
)

type VoteMapper struct{}

func NewVoteMapper() *VoteMapper {
	return &VoteMapper{}
}

func (m *VoteMapper) RequestToDomain(req *dto.VoteRequestDTO, session string) (domain.VoteEntry, error) {
	entry := domain.VoteEntry{
		ID:        uuid.New().String(),
		Session:   session,
		Value:     value.VoteValue(req.Value),
		Timestamp: req.TimeStamp,
	}

	return entry, entry.Valid()
}

func (m *VoteMapper) DomainToResponse(votes []domain.VoteEntry) []dto.VoteResponseDTO {
	out := make([]dto.VoteResponseDTO, 0, len(votes))

	for _, ve := range votes {
		out = append(out, dto.VoteResponseDTO{
			Value:     int(ve.Value),
			Timestamp: ve.Timestamp,
		})
	}

	return out
}
