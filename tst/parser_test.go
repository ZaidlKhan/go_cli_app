package tst

import (
	"main/internal/parser"
	"main/internal/types"
	"os"
	"testing"

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
	_, proc := setupTest(t)
	expected := "http://127.0.0.1:3000"

	proc.URL = ""
	parser.ResolveURL(proc)

	assert.Equal(t, expected, proc.URL)
}
