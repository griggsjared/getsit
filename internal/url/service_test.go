package url_test

import (
	"context"
	"testing"

	"github.com/griggsjared/getsit/internal/url"
	"github.com/griggsjared/getsit/internal/url/entity"
	"github.com/griggsjared/getsit/internal/url/repository"
)

func TestService_SaveUrl(t *testing.T) {

	ctx := context.Background()
	r := repository.NewMemUrlEntryRepository()
	s := url.NewService(r)

	//save an existing url
	_, err := s.SaveUrl(ctx, &url.SaveUrlInput{
		Url: "https://exists.com",
	})
	if err != nil {
		t.Errorf("SaveUrl() error = %v", err)
	}

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid url",
			url:     "https://new.com",
			wantErr: false,
		},
		{
			name:    "valid url but already exists",
			url:     "https://exists.com",
			wantErr: true,
		},
		{
			name:    "invalid url",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "empty url",
			url:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.SaveUrl(ctx, &url.SaveUrlInput{
				Url: tt.url,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveUrl() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestService_GetUrlByToken(t *testing.T) {

	ctx := context.Background()
	r := repository.NewMemUrlEntryRepository()
	s := url.NewService(r)

	entry, err := s.SaveUrl(ctx, &url.SaveUrlInput{
		Url: "https://example.com",
	})
	if err != nil {
		t.Errorf("SaveUrl() error = %v", err)
	}

  token, err := entity.NewUrlToken()
  if err != nil {
    t.Errorf("NewUrlToken() error = %v", err)
  }

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token and found",
			token:   entry.Token.String(),
			wantErr: false,
		},
		{
			name:    "valid token but not found",
			token:   token.String(),
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "abc123",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetUrlByToken(ctx, &url.GetUrlByTokenInput{
				Token: tt.token,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrlByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}

}

func TestService_GetUrlByUrl(t *testing.T) {

	ctx := context.Background()
	r := repository.NewMemUrlEntryRepository()
	s := url.NewService(r)

	entry, err := s.SaveUrl(ctx, &url.SaveUrlInput{
		Url: "https://example.com",
	})
	if err != nil {
		t.Errorf("SaveUrl() error = %v", err)
	}

	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "valid short url and found",
			url:     entry.Url.String(),
			wantErr: false,
		},
		{
			name:    "valid short url but not found",
			url:     "https://notfound.com",
			wantErr: true,
		},
		{
			name:    "invalid short url",
			url:     "example.com",
			wantErr: true,
		},
		{
			name:    "empty short url",
			url:     "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := s.GetUrlByUrl(ctx, &url.GetUrlByUrlInput{
				Url: tt.url,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUrlByUrll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestService_VisitUrlByToken(t *testing.T) {

	ctx := context.Background()
	r := repository.NewMemUrlEntryRepository()
	s := url.NewService(r)

	entry, err := s.SaveUrl(ctx, &url.SaveUrlInput{
		Url: "https://example.com",
	})
	if err != nil {
		t.Errorf("SaveUrl() error = %v", err)
	}

  token, err := entity.NewUrlToken()
  if err != nil {
    t.Errorf("NewUrlToken() error = %v", err)
  }

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "valid token and found",
			token:   entry.Token.String(),
			wantErr: false,
		},
		{
			name:    "valid token but not found",
			token:   token.String(),
			wantErr: true,
		},
		{
			name:    "invalid token",
			token:   "abc123",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.VisitUrlByToken(ctx, &url.VisitUrlByTokenInput{
				Token: tt.token,
			})
			if (err != nil) != tt.wantErr {
				t.Errorf("VisitUrlByToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
