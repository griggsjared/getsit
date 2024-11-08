package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed public
var assets embed.FS

// AssetsFS returns the embedded filesystem for the assets and strips the public prefix from the path.
func AssetsFS() http.FileSystem {
	fileS, err := fs.Sub(assets, "public")
	if err != nil {
		log.Fatal(err)
	}
	return http.FS(fileS)
}
