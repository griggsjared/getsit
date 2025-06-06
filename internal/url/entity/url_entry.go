package entity

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"regexp"
)

const (
	tokenBytes        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_"
	tokenLength uint8 = 8
)

// UrlToken is a short string that will be used to access the long url.
// This token is 8 characters long and can be and upper or lower case letter or a number
type UrlToken string

// NewUrlToken will generate a new url token that is 8 characters long.
func NewUrlToken() (UrlToken, error) {
	var token UrlToken
	tokenBytesLen := big.NewInt(int64(len(tokenBytes)))
	for i := 0; i < int(tokenLength); i++ {
		randomIdx, err := rand.Int(rand.Reader, tokenBytesLen)
		if err != nil {
			return "", err
		}
		token += UrlToken(tokenBytes[randomIdx.Int64()])
	}

	return token, nil
}

// Validate will check if the token is valid
func (t UrlToken) Validate() error {
	//use regex to check if the token is valid
	mustCompile := regexp.MustCompile(`^[a-zA-Z0-9_]{8}$`)
	if !mustCompile.MatchString(t.String()) {
		return fmt.Errorf("token is not valid")
	}
	return nil
}

// String will return the string representation of the token
func (t UrlToken) String() string {
	return string(t)
}

// Url is the long url the is associated with the token
type Url string

// Validate will check if the url is valid
func (u Url) Validate() error {

	// Check if the url is a valid url
	pu, err := url.Parse(u.String())
	if err != nil {
		return fmt.Errorf("url is not valid")
	}

	if pu.Scheme == "" || pu.Host == "" {
		return fmt.Errorf("url is not valid")
	}

	return nil
}

// String will return the string representation of the url
func (u Url) String() string {
	return string(u)
}

// UrlEntry is the domain entity that will store the long url, token, and the number of times the url has been visited
type UrlEntry struct {
	Url        Url      // The long url
	Token      UrlToken // The token is a short string that will be used to access the long url
	VisitCount int      // The number of times the url has been visited
}

// NewUrlEntry will create a new url entry from primitive types
func NewUrlEntry(url string, token string, visitCount int) *UrlEntry {
	return &UrlEntry{
		Url:        Url(url),
		Token:      UrlToken(token),
		VisitCount: visitCount,
	}
}
