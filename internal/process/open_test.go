package process

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestBrowserURLPreservesHTTPWhenEndpointSupportsTLS(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	inputURL := strings.Replace(server.URL, "https://", "http://", 1)
	openedURL, err := browserURL(inputURL)
	if err != nil {
		t.Fatalf("browserURL returned error: %v", err)
	}
	if openedURL.Scheme != "http" {
		t.Errorf("scheme = %q, want http", openedURL.Scheme)
	}
}

func TestBrowserURLKeepsHTTPWhenEndpointDoesNotSupportTLS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	openedURL, err := browserURL(server.URL)
	if err != nil {
		t.Fatalf("browserURL returned error: %v", err)
	}
	if openedURL.Scheme != "http" {
		t.Errorf("scheme = %q, want http", openedURL.Scheme)
	}
}

func TestBrowserURLUsesIPv4LoopbackForLocalhost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	localURL := strings.Replace(server.URL, "127.0.0.1", "localhost", 1)
	openedURL, err := browserURL(localURL)
	if err != nil {
		t.Fatalf("browserURL returned error: %v", err)
	}
	if got, want := openedURL.Hostname(), "127.0.0.1"; got != want {
		t.Errorf("hostname = %q, want %q", got, want)
	}
}

func TestBrowserURLAddsHTTPWhenSchemeIsMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	rawURL := strings.TrimPrefix(server.URL, "http://")
	openedURL, err := browserURL(rawURL)
	if err != nil {
		t.Fatalf("browserURL returned error: %v", err)
	}
	if got, want := openedURL.Scheme, "http"; got != want {
		t.Errorf("scheme = %q, want %q", got, want)
	}
	if got, want := openedURL.Host, strings.TrimPrefix(server.URL, "http://"); got != want {
		t.Errorf("host = %q, want %q", got, want)
	}
}
