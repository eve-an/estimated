package dto

type EventResponseDTO struct {
	Votes map[string][]VoteResponseDTO `json:"votes"`
}
