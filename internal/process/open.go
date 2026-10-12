package process

import (
	"fmt"
	"net"
	"net/url"
	"os/exec"
	"strings"
)

func OpenURL(rawURL string) error {
	parsedURL, err := browserURL(rawURL)
	if err != nil {
		return err
	}

	if err := exec.Command("open", parsedURL.String()).Run(); err != nil {
		return fmt.Errorf("launch browser: %w", err)
	}
	return nil
}

func browserURL(rawURL string) (*url.URL, error) {
	if !strings.Contains(rawURL, "://") {
		rawURL = "http://" + rawURL
	}

	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, fmt.Errorf("parse URL: %w", err)
	}
	if parsedURL.Host == "" || (parsedURL.Scheme != "http" && parsedURL.Scheme != "https") {
		return nil, fmt.Errorf("unsupported URL: %q", rawURL)
	}

	if strings.EqualFold(parsedURL.Hostname(), "localhost") {
		if port := parsedURL.Port(); port != "" {
			parsedURL.Host = net.JoinHostPort("127.0.0.1", port)
		} else {
			parsedURL.Host = "127.0.0.1"
		}
	}
	return parsedURL, nil
}
