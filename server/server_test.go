package server_test

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/brighteyed/http-server/config"
	"github.com/brighteyed/http-server/server"
)

func TestFileSystemHandle(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("handle existing file", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata/html"}}

		request, err := http.NewRequest(http.MethodGet, "/test/welcome.html", nil)
		if err != nil {
			t.Fatal("Cannot create request")
		}
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 200)

		expectedBody, err := os.ReadFile("testdata/html/welcome.html")
		if err != nil {
			t.Fatal("Error reading test file")
		}

		if !bytes.Equal(expectedBody, response.Body.Bytes()) {
			t.Error("Got unexpected response body")
		}
	})

	t.Run("handle not existing file", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata/html"}}

		request, err := http.NewRequest(http.MethodGet, "/test/credits.html", nil)
		if err != nil {
			t.Fatal("Cannot create request")
		}
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 404)
	})
}

func assertStatusCode(t *testing.T, got, expected int) {
	t.Helper()

	if got != expected {
		t.Fatalf("Expected Status code %d, but got %d", expected, got)
	}
}
