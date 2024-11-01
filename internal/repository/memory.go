package repository

import (
	"context"
	"fmt"

	"github.com/griggsjared/getsit/internal/entity"
)

// memEntriesTokenMap is a map that will store the url entry with the token as the key
type memEntriesTokenMap map[entity.UrlToken]*entity.UrlEntry

// memEntriesUrlMap is a map that will store the url entry with the url as the key
type memEntriesUrlMap map[entity.Url]*entity.UrlEntry

// MemUrlEntryStore is a in memory store that will store the url entries
type MemUrlEntryStore struct {
	entriesToken memEntriesTokenMap //key is the token and value is the url entry for a fast lookup ( O(1) )
	entriesUrl   memEntriesUrlMap   //key is the url and value is the url entry for a fast lookup ( O(1) )
}

// NewMemUrlEntryStore will create a new in memory store
func NewMemUrlEntryStore() *MemUrlEntryStore {
	return &MemUrlEntryStore{
		entriesToken: make(memEntriesTokenMap),
		entriesUrl:   make(memEntriesUrlMap),
	}
}

// Save will save the url entry to the store
func (s *MemUrlEntryStore) SaveUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

	var entry *entity.UrlEntry

	if e, ok := s.entriesUrl[url]; ok {
		entry = e
	} else {
		var token entity.UrlToken
		for {
			token = entity.NewUrlToken()
			if _, ok := s.entriesToken[token]; !ok {
				break
			}
		}
		entry = &entity.UrlEntry{
			Url:        url,
			Token:      token,
			VisitCount: 0,
		}
	}

	s.entriesToken[entry.Token] = entry
	s.entriesUrl[entry.Url] = entry

	return entry, nil
}

// SaveVisit will increment the number of times the url has been visited
func (s *MemUrlEntryStore) SaveVisit(ctx context.Context, token entity.UrlToken) error {
	if e, ok := s.entriesToken[token]; ok {
		e.VisitCount++
		return nil
	}
	return fmt.Errorf("entry not found")
}

// GetFromToken will return the url entry for the given token
func (s *MemUrlEntryStore) GetFromToken(ctx context.Context, token entity.UrlToken) (*entity.UrlEntry, error) {
	if e, ok := s.entriesToken[token]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("entry not found")
}

// GetFromUrl will return the url entry for the given url
func (s *MemUrlEntryStore) GetFromUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {
	if e, ok := s.entriesUrl[url]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("entry not found")
}
