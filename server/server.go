package server

import (
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/brighteyed/http-server/config"
)

// Handler handles requests for the server
type Handler struct {
	http.Handler
}

// NewHandler returns a Handler created by the slice of
// locations
func NewHandler(locations []config.Location) *Handler {
	routes := make(map[string]bool)
	router := http.NewServeMux()

	for i := 0; i < len(locations); i++ {
		path := locations[i].Path
		root := locations[i].Root

		if _, exists := routes[path]; exists {
			log.Printf("Route %q is already in use\n", path)
			continue
		}

		if strings.ToLower(filepath.Ext(root)) != ".zip" {
			router.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(root))))
		} else {
			h := NewZipHandler(root)
			router.Handle(path, http.StripPrefix(path, http.HandlerFunc(h.GetFile)))
		}

		routes[path] = true
		log.Printf("Serving %q as %q\n", root, path)
	}

	return &Handler{
		Handler: router,
	}
}
