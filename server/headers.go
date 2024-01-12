package server

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Header struct {
	Name  string
	Value string
}

type HeaderList []Header

func (h *HeaderList) Set(value string) error {
	parts := strings.SplitN(value, ":", 2)
	if len(parts) != 2 {
		return errors.New("header must be in the format 'name: value'")
	}
	*h = append(*h, Header{Name: strings.TrimSpace(parts[0]), Value: strings.TrimSpace(parts[1])})

	return nil
}

func (h *HeaderList) String() string {
	return fmt.Sprintf("%v", *h)
}

func AddHeaders(headers HeaderList, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, h := range headers {
			w.Header().Set(h.Name, h.Value)
		}
		next.ServeHTTP(w, r)
	})
}
