package parser_test

import (
	"os"
	"testing"

	"main/internal/parser"
	"main/internal/types"

	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (string, *types.Process) {
	t.Helper()

	tmpDir, err := os.MkdirTemp("", "pct-test-*")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	proc := &types.Process{
		PID:         1234,
		ProcessName: "node",
		Address:     "127.0.0.1:3000",
		Port:        3000,
		URL:         "http://127.0.0.1:3000",
		GitBranch:   "main",
	}

	return tmpDir, proc
}

func TestResolveURL(t *testing.T) {
	tests := []struct {
		name     string
		address  string
		port     uint16
		expected string
	}{
		{
			name:     "standard localhost ipv4",
			address:  "127.0.0.1:3000",
			port:     3000,
			expected: "http://127.0.0.1:3000",
		},
		{
			name:     "ipv6 loopback address",
			address:  "[::1]:8080",
			port:     8080,
			expected: "http://[::1]:8080",
		},
		{
			name:     "zero port returns empty url",
			address:  "127.0.0.1:0",
			port:     0,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, proc := setupTest(t)
			proc.Address = tt.address
			proc.Port = tt.port
			proc.URL = ""

			parser.ResolveURL(proc)

			assert.Equal(t, tt.expected, proc.URL)
		})
	}
}

func TestInspectPath(t *testing.T) {
	tests := []struct {
		name     string
		cmdArgs  string
		expected bool
	}{
		{
			name:     "valid application path",
			cmdArgs:  "node app.js",
			expected: true,
		},
		{
			name:     "invalid system prefix",
			cmdArgs:  "/System/Library/CoreServices/Finder.app",
			expected: false,
		},
		{
			name:     "empty path string",
			cmdArgs:  "",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid := parser.InspectPath(tt.cmdArgs)
			assert.Equal(t, tt.expected, isValid)
		})
	}
}

func TestParseOutput(t *testing.T) {
	mockInput := "p1234\ncnode\nn127.0.0.1:3000\n"

	processes := parser.ParseOutput(mockInput)

	if len(processes) > 0 {
		assert.Contains(t, processes, 1234)
		assert.Equal(t, uint16(3000), processes[1234].Port)
		assert.Equal(t, "127.0.0.1:3000", processes[1234].Address)
	}
}
