package server

import (
	"log"
	"net/http"

	"github.com/brighteyed/http-server/config"
)

// Handler handles requests for the server
type Handler struct {
	http.Handler
}

// NewHandler returns a Handler created by the slice of
// locations
func NewHandler(locations []config.Location) *Handler {
	router := http.NewServeMux()

	for i := 0; i < len(locations); i++ {
		path := locations[i].Path
		root := locations[i].Root

		router.Handle(path, http.StripPrefix(path, http.FileServer(http.Dir(root))))
		log.Printf("Serving %q as %q\n", root, path)
	}

	return &Handler{
		Handler: router,
	}
}
