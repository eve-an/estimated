package model

import "time"

type VoteEntry struct {
	Value     int       `json:"value"`
	TimeStamp time.Time `json:"timestamp"`
}
