package entity_test

import (
	"testing"

	"github.com/griggsjared/getsit/internal/entity"
)

func TestUrlToken_NewUrlToken(t *testing.T) {
	token := entity.NewUrlToken()
	if err := token.Validate(); err != nil {
		t.Errorf("NewUrlToken() = %v", err)
	}
}

func TestUrlToken_Validate(t *testing.T) {
	tests := []struct {
		name    string
		token   entity.UrlToken
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   entity.UrlToken("abc12345"),
			wantErr: false,
		},
		{
			name:    "invalid token",
			token:   entity.UrlToken("abc123"),
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   entity.UrlToken(""),
			wantErr: true,
		},
		{
			name:    "token with special characters",
			token:   entity.UrlToken("abc12345!"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.token.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("UrlToken.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUrlToken_String(t *testing.T) {
	tests := []struct {
		name string
		t    entity.UrlToken
		want string
	}{
		{
			name: "valid token",
			t:    entity.UrlToken("abc12345"),
			want: "abc12345",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.t.String(); got != tt.want {
				t.Errorf("UrlToken.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUrl_Validate(t *testing.T) {
	tests := []struct {
		name    string
		u       entity.Url
		wantErr bool
	}{
		{
			name:    "valid https url",
			u:       entity.Url("https://example.com"),
			wantErr: false,
		},
		{
			name:    "valid http url",
			u:       entity.Url("https://example.com"),
			wantErr: false,
		},
		{
			name:    "invalid url with no scheme",
			u:       entity.Url("example.com"),
			wantErr: true,
		},
		{
			name:    "empty url",
			u:       entity.Url(""),
			wantErr: true,
		},
		{
			name:    "invalid url with no host",
			u:       entity.Url("https://"),
			wantErr: true,
		},
		{
			name:    "invalid url with no scheme or host",
			u:       entity.Url("example"),
			wantErr: true,
		},
		{
			name:    "invalid url a :",
			u:       entity.Url(":"),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Url.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUrl_String(t *testing.T) {
	tests := []struct {
		name string
		u    entity.Url
		want string
	}{
		{
			name: "valid url",
			u:    entity.Url("https://example.com"),
			want: "https://example.com",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.String(); got != tt.want {
				t.Errorf("Url.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewUrlEntry(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		token      string
		visitCount int
	}{
		{
			name:       "valid url entry",
			url:        "https://example.com",
			token:      "abc12345",
			visitCount: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := entity.NewUrlEntry(tt.url, tt.token, tt.visitCount)
			if got.Url.String() != tt.url {
				t.Errorf("NewUrlEntry() = %v, want %v", got.Url.String(), tt.url)
			}
			if got.Token.String() != tt.token {
				t.Errorf("NewUrlEntry() = %v, want %v", got.Token.String(), tt.token)
			}
			if got.VisitCount != tt.visitCount {
				t.Errorf("NewUrlEntry() = %v, want %v", got.VisitCount, tt.visitCount)
			}
		})
	}
}
