package endpoint

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

type staticSite struct {
	h http.Handler
}

//go:embed static
var staticFiles embed.FS

func NewStaticSite() http.Handler {
	fsys, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("failed to get sub fs: %v", err)
	}

	return &staticSite{
		h: http.FileServer(http.FS(fsys)),
	}
}

func (s staticSite) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.h.ServeHTTP(w, r)
}
