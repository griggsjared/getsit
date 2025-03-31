package repository

import (
	"context"
	"fmt"

	"github.com/griggsjared/getsit/internal/url/entity"
)

// memEntriesTokenMap is a map that will repository the url entry with the token as the key
type memEntriesTokenMap map[entity.UrlToken]*entity.UrlEntry

// memEntriesUrlMap is a map that will repository the url entry with the url as the key
type memEntriesUrlMap map[entity.Url]*entity.UrlEntry

// MemUrlEntryRepository is a in memory repository that will repository the url entries
type MemUrlEntryRepository struct {
	entriesToken memEntriesTokenMap //key is the token and value is the url entry for a fast lookup ( O(1) )
	entriesUrl   memEntriesUrlMap   //key is the url and value is the url entry for a fast lookup ( O(1) )
}

// NewMemUrlEntryRepository will create a new in memory repository
func NewMemUrlEntryRepository() *MemUrlEntryRepository {
	return &MemUrlEntryRepository{
		entriesToken: make(memEntriesTokenMap),
		entriesUrl:   make(memEntriesUrlMap),
	}
}

// Save will save the url entry to the repository
func (s *MemUrlEntryRepository) SaveUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {

	var entry *entity.UrlEntry

	//if the url already exists escape with an error
	if _, ok := s.entriesUrl[url]; ok {
		return nil, fmt.Errorf("entry already exists")
	}

	var token entity.UrlToken
	var err error
	for {
		token, err = entity.NewUrlToken()
		if err != nil {
			return nil, err
		}
		if _, ok := s.entriesToken[token]; !ok {
			break
		}
	}
	entry = &entity.UrlEntry{
		Url:        url,
		Token:      token,
		VisitCount: 0,
	}

	s.entriesToken[entry.Token] = entry
	s.entriesUrl[entry.Url] = entry

	return entry, nil
}

// SaveVisit will increment the number of times the url has been visited
func (s *MemUrlEntryRepository) SaveVisit(ctx context.Context, token entity.UrlToken) error {
	if e, ok := s.entriesToken[token]; ok {
		e.VisitCount++
		return nil
	}
	return fmt.Errorf("entry not found")
}

// GetFromToken will return the url entry for the given token
func (s *MemUrlEntryRepository) GetFromToken(ctx context.Context, token entity.UrlToken) (*entity.UrlEntry, error) {
	if e, ok := s.entriesToken[token]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("entry not found")
}

// GetFromUrl will return the url entry for the given url
func (s *MemUrlEntryRepository) GetFromUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error) {
	if e, ok := s.entriesUrl[url]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("entry not found")
}
