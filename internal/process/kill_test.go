package process

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKillProcess_Success(t *testing.T) {
	cmd := exec.Command("sleep", "10")
	err := cmd.Start()
	assert.NoError(t, err, "failed to start helper process")

	pid := cmd.Process.Pid

	err = KillProcess(pid)
	assert.NoError(t, err, "expected no error killing process %d", pid)

	err = cmd.Wait()
	assert.Error(t, err, "expected helper process to be terminated, but it exited normally")
}

func TestKillProcess_Error(t *testing.T) {
	invalidPID := -99999

	err := KillProcess(invalidPID)
	assert.Error(t, err, "expected an error when attempting to kill invalid PID %d", invalidPID)
}
