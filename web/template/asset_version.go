package template

import (
	"crypto/rand"
	"fmt"
)

// assetVersion is package level variable that holds a random string for cache busting
var assetVersion string

// init generates a random string for the assetVersion
func init() {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	assetVersion = fmt.Sprintf("%x", b)
}
