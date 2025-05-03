package db

import "github.com/eve-an/estimated/internal/model"

type VoteEntryStore interface {
	// Add an entry to the store
	Add(key string, vote model.VoteEntry) error
	// Get the entry with the assosiated key from the store
	Get(key string) ([]model.VoteEntry, error)
	// List all entries from the store
	List() ([]model.VoteEntry, error)
	// Clear all entries from the store
	Clear() (int, error)
}
