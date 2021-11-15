package server

import (
	"archive/zip"
	"errors"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
)

// ZipHandler implements http.Handler interface to serve contents of
// a zip file
type ZipHandler struct {
	file string
}

// GetFile acts as http.HandleFunc to serve contents of a zip file
func (z *ZipHandler) GetFile(w http.ResponseWriter, r *http.Request) {
	const defaultFile = "index.html"

	rc, err := zip.OpenReader(z.file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rc.Close()

	path := strings.TrimSuffix(r.URL.Path, "/")
	if path == "" {
		path = defaultFile
	}

	file, err := rc.Open(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			http.Error(w, "404 page not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if info.IsDir() {
		path += "/" + defaultFile

		file, err = rc.Open(path)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				http.Error(w, "404 page not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		defer file.Close()
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ctype := mime.TypeByExtension(filepath.Ext(r.URL.Path))
	if ctype == "" {
		ctype = http.DetectContentType(bytes)
	}
	w.Header().Set("Content-Type", ctype)
	if _, err := w.Write(bytes); err != nil {
		log.Printf("Error writing response, %v", err)
	}
}

func NewZipHandler(zipfile string) *ZipHandler {
	return &ZipHandler{zipfile}
}
