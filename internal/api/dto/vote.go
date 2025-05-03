package dto

import (
	"time"
)

type VoteRequestDTO struct {
	Value     int       `json:"value"`
	TimeStamp time.Time `json:"timestamp"`
}

type VoteResponseDTO struct {
	Value     int       `json:"value"`
	Timestamp time.Time `json:"timestamp"`
}
