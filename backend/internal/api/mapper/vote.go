package mapper

import (
	"log/slog"

	"github.com/eve-an/estimated/internal/api/dto"
	"github.com/eve-an/estimated/internal/domain"
	"github.com/eve-an/estimated/internal/value"
	"github.com/google/uuid"
)

type VoteMapper struct {
	logger *slog.Logger
}

func NewVoteMapper(logger *slog.Logger) *VoteMapper {
	return &VoteMapper{logger: logger}
}

func (m *VoteMapper) RequestToDomain(req *dto.VoteRequestDTO, session string, name string) (domain.VoteEntry, error) {
	return domain.NewVoteEntry(
		uuid.New().String(),
		session,
		name,
		value.VoteValue(req.Value),
		req.TimeStamp,
	)
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
