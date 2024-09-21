package repository

import "github.com/griggsjared/getsit/internal/entity"

type UrlEntryRepository interface {
	// Save will url entry to the store
	Save(url entity.Url) (entry *entity.UrlEntry, new bool, err error)
	// SaveVisit will increment the number of times the url has been visited
	SaveVisit(token string) error
	// Get will return the url entry for the given token
	GetFromToken(token string) (*entity.UrlEntry, error)
	// Get will return the url entry for the given url
	GetFromUrl(url string) (*entity.UrlEntry, error)
}
