package repository

import (
	"fmt"

	"github.com/griggsjared/getsit/internal/entity"
)

// entriesTokenMap is a map that will store the url entry with the token as the key
type entriesTokenMap map[entity.UrlToken]*entity.UrlEntry

// entriesUrlMap is a map that will store the url entry with the url as the key
type entriesUrlMap map[entity.Url]*entity.UrlEntry

// InMemoryUrlEntryStore is a in memory store that will store the url entries
type InMemoryUrlEntryStore struct {
	entriesToken entriesTokenMap //key is the token and value is the url entry for a fast lookup ( O(1) )
	entriesUrl   entriesUrlMap   //key is the url and value is the url entry for a fast lookup ( O(1) )
}

// NewInMemoryUrlEntryStore will create a new in memory store
func NewInMemoryUrlEntryStore() *InMemoryUrlEntryStore {
	return &InMemoryUrlEntryStore{
		entriesToken: make(entriesTokenMap),
		entriesUrl:   make(entriesUrlMap),
	}
}

// Save will save the url entry to the store
func (s *InMemoryUrlEntryStore) Save(url entity.Url) (urlEntry *entity.UrlEntry, new bool, err error) {

	var entry *entity.UrlEntry

	// Check if the url entry already exists
	// if it does we just grab the entry
	// if it does not we create a new entry
	var exists bool
	if e, ok := s.entriesUrl[url]; ok {
		exists = true
		entry = e
	} else {
		exists = false
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

	return entry, !exists, nil
}

// SaveVisit will increment the number of times the url has been visited
func (s *InMemoryUrlEntryStore) SaveVisit(token string) error {
	if e, ok := s.entriesToken[entity.UrlToken(token)]; ok {
		e.VisitCount++
		return nil
	}
	return fmt.Errorf("entry not found")
}

// GetFromToken will return the url entry for the given token
func (s *InMemoryUrlEntryStore) GetFromToken(token string) (*entity.UrlEntry, error) {
	if e, ok := s.entriesToken[entity.UrlToken(token)]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("entry not found")
}

// GetFromUrl will return the url entry for the given url
func (s *InMemoryUrlEntryStore) GetFromUrl(url string) (*entity.UrlEntry, error) {
	if e, ok := s.entriesUrl[entity.Url(url)]; ok {
		return e, nil
	}
	return nil, fmt.Errorf("entry not found")
}
