package process

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBrowserURLPreservesHTTPWhenEndpointSupportsTLS(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	inputURL := strings.Replace(server.URL, "https://", "http://", 1)
	openedURL, err := browserURL(inputURL)
	require.NoError(t, err, "browserURL returned error")
	assert.Equal(t, "http", openedURL.Scheme)
}

func TestBrowserURLKeepsHTTPWhenEndpointDoesNotSupportTLS(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	openedURL, err := browserURL(server.URL)
	require.NoError(t, err, "browserURL returned error")
	assert.Equal(t, "http", openedURL.Scheme)
}

func TestBrowserURLUsesIPv4LoopbackForLocalhost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	localURL := strings.Replace(server.URL, "127.0.0.1", "localhost", 1)
	openedURL, err := browserURL(localURL)
	require.NoError(t, err, "browserURL returned error")
	assert.Equal(t, "127.0.0.1", openedURL.Hostname())
}

func TestBrowserURLAddsHTTPWhenSchemeIsMissing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	defer server.Close()

	rawURL := strings.TrimPrefix(server.URL, "http://")
	openedURL, err := browserURL(rawURL)
	require.NoError(t, err, "browserURL returned error")
	assert.Equal(t, "http", openedURL.Scheme)
	assert.Equal(t, strings.TrimPrefix(server.URL, "http://"), openedURL.Host)
}
