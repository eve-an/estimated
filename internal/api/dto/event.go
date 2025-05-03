package dto

import (
	"github.com/eve-an/estimated/internal/domain"
)

type EventResponseDTO struct {
	Name  string             `json:"name"`
	Votes []domain.VoteEntry `json:"votes"`
}
