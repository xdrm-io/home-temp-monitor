package endpoint

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

//go:embed static/*
var staticFS embed.FS

// StaticWeb is a static web server from the embed static files
type StaticWeb http.Handler

// NewStaticWeb returns a new StaticWeb instance
func NewStaticWeb() StaticWeb {
	fsys, err := fs.Sub(staticFS, "static")
	if err != nil {
		log.Fatalf("cannot load static files: %v", err)
	}
	return http.FileServer(http.FS(fsys))
}
