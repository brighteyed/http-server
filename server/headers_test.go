package server_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brighteyed/http-server/server"
)

func TestAddHeaders(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for name, values := range r.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}
	})

	middleware := server.AddHeaders(server.HeaderList{
		{Name: "Content-Type", Value: "application/json"},
		{Name: "Cache-Control", Value: "no-cache"},
	}, handler)

	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set("User-Agent", "TestAgent")

	response := httptest.NewRecorder()
	middleware.ServeHTTP(response, request)

	expectedHeaders := map[string][]string{
		"Content-Type":  {"application/json"},
		"Cache-Control": {"no-cache"},
		"User-Agent":    {"TestAgent"},
	}

	for name, expectedValues := range expectedHeaders {
		actualValues := response.Header()[name]
		if !equalStringSlices(actualValues, expectedValues) {
			t.Errorf("Expected header %s to have values %v, but got %v", name, expectedValues, actualValues)
		}
	}
}

func equalStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
