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
			{Path: "/test/", Root: "testdata"}}

		request := mustGetRequest(t, "/test/welcome.html")
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 200)

		assertContentType(t, response.Header().Get("Content-Type"), "text/html; charset=utf-8")
		assertContent(t, response.Body, "testdata/welcome.html")
	})

	t.Run("handle not existing file", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata"}}

		request := mustGetRequest(t, "/test/credits.html")
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 404)
	})
}

func TestZipHandle(t *testing.T) {
	log.SetOutput(ioutil.Discard)

	t.Run("handle existing file", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata/archive.zip"}}

		request := mustGetRequest(t, "/test/welcome.html")
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 200)

		assertContentType(t, response.Header().Get("Content-Type"), "text/html; charset=utf-8")
		assertContent(t, response.Body, "testdata/welcome.html")
	})

	t.Run("handle not existing file", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata/archive.zip"}}

		request := mustGetRequest(t, "/test/credits.html")
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 404)
	})

	t.Run("handle root", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata/archive.zip"}}

		request := mustGetRequest(t, "/test/")
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 404)
	})

	t.Run("handle directory without index.html", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata/archive.zip"}}

		request := mustGetRequest(t, "/test/empty")
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 404)
	})

	t.Run("handle directory with index.html", func(t *testing.T) {
		locations := []config.Location{
			{Path: "/test/", Root: "testdata/archive.zip"}}

		request := mustGetRequest(t, "/test/subdir/")
		response := httptest.NewRecorder()

		handler := server.NewHandler(locations)
		handler.ServeHTTP(response, request)

		gotStatusCode := response.Result().StatusCode
		assertStatusCode(t, gotStatusCode, 200)
		assertContent(t, response.Body, "testdata/index.html")
	})
}

func mustGetRequest(t *testing.T, url string) *http.Request {
	t.Helper()

	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		t.Fatal("Cannot create request")
	}
	return request
}

func assertStatusCode(t *testing.T, got, expected int) {
	t.Helper()

	if got != expected {
		t.Fatalf("Expected Status code %d, but got %d", expected, got)
	}
}

func assertContentType(t *testing.T, got, expected string) {
	t.Helper()

	if got != expected {
		t.Errorf("Expected %q, got %q", expected, got)
	}
}

func assertContent(t *testing.T, buffer *bytes.Buffer, file string) {
	t.Helper()

	expectedBody, err := os.ReadFile(file)
	if err != nil {
		t.Fatal("Error reading test file")
	}

	if !bytes.Equal(expectedBody, buffer.Bytes()) {
		t.Error("Got unexpected response body")
	}
}
