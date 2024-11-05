package service

import (
	"context"
	"errors"

	"github.com/griggsjared/getsit/internal/entity"
)

// ErrValidation is a generic validation error that can be returned when input validation fails
var ErrValidation = errors.New("validation error")

// withValidationErrors is a struct that can be embedded into the various input structs to hold validation errors
type withValidationErrors struct {
	ValidationErrors map[string]string
}

// UrlEntryRepository is the interface that defines the method that the service will use to interact with the repository
type UrlEntryRepository interface {
	// Save will url entry to the store
	SaveUrl(ctx context.Context, url entity.Url) (entry *entity.UrlEntry, err error)
	// SaveVisit will increment the number of times the url has been visited
	SaveVisit(ctx context.Context, token entity.UrlToken) error
	// GetFromToken will get the url entry from the token
	GetFromToken(ctx context.Context, token entity.UrlToken) (*entity.UrlEntry, error)
	// GetFromUrl will get the url entry from the url
	GetFromUrl(ctx context.Context, url entity.Url) (*entity.UrlEntry, error)
}

type Service struct {
	repo UrlEntryRepository
}

// New will create a new service
func New(repo UrlEntryRepository) *Service {
	return &Service{
		repo: repo,
	}
}

// SaveUrlInput is the input struct for the SaveUrl method
type SaveUrlInput struct {
	withValidationErrors
	Url string
}

// SaveUrl will validate the url string and save it to the store
func (s *Service) SaveUrl(ctx context.Context, input *SaveUrlInput) (*entity.UrlEntry, error) {

	input.ValidationErrors = make(map[string]string)

	// Validate the url
	urlEntry := entity.Url(input.Url)
	if err := urlEntry.Validate(); err != nil {
		input.ValidationErrors["url"] = err.Error()
		return nil, ErrValidation
	}

	// Save the url
	entry, err := s.repo.SaveUrl(ctx, urlEntry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

// GetUrlInput is the input struct for the GetUrl method
type GetUrlByTokenInput struct {
	withValidationErrors
	Token string
}

// GetUrl will get the url entry from the url string
func (s *Service) GetUrlByToken(ctx context.Context, input *GetUrlByTokenInput) (*entity.UrlEntry, error) {

	input.ValidationErrors = make(map[string]string)

	// Validate the token
	token := entity.UrlToken(input.Token)
	if err := token.Validate(); err != nil {
		input.ValidationErrors["token"] = err.Error()
		return nil, ErrValidation
	}

	// Get the url entry
	entry, err := s.repo.GetFromToken(ctx, token)
	if err != nil {
		return nil, errors.New("failed to get url")
	}

	return entry, nil
}

// GetUrlInput is the input struct for the GetUrl method
type GetUrlByUrlInput struct {
	withValidationErrors
	Url string
}

// GetUrl will get the url entry from the url string
func (s *Service) GetUrlByUrl(ctx context.Context, input *GetUrlByUrlInput) (*entity.UrlEntry, error) {

	input.ValidationErrors = make(map[string]string)

	// Validate the url
	urlEntry := entity.Url(input.Url)
	if err := urlEntry.Validate(); err != nil {
		input.ValidationErrors["url"] = err.Error()
		return nil, ErrValidation
	}

	// Get the url entry
	entry, err := s.repo.GetFromUrl(ctx, urlEntry)
	if err != nil {
		return nil, errors.New("failed to get url")
	}

	return entry, nil
}

// VisitUrlInput is the input struct for the VisitUrl method
type VisitUrlByTokenInput struct {
	withValidationErrors
	Token string
}

// VisitUrl will increment the number of times the url has been visited
func (s *Service) VisitUrlByToken(ctx context.Context, input *VisitUrlByTokenInput) error {

	input.ValidationErrors = make(map[string]string)

	// Validate the token
	urlToken := entity.UrlToken(input.Token)
	if err := urlToken.Validate(); err != nil {
		input.ValidationErrors["token"] = err.Error()
		return ErrValidation
	}

	// Save the visit
	err := s.repo.SaveVisit(ctx, urlToken)
	if err != nil {
		return err
	}

	return nil
}
