// Public Domain (-) 2018-present, The Peerbase Authors.
// See the Peerbase UNLICENSE file for details.

// +build !windows

package process

import (
	"os"
	"os/exec"
	"syscall"
	"testing"
)

func TestReapOrphans(t *testing.T) {
	// First, pretend to be PID 1.
	osGetpid = func() int {
		return 1
	}
	testMode = true
	cmd := exec.Command("sleep", "100")
	if err := cmd.Start(); err != nil {
		t.Fatalf("Unexpected error when trying to run `sleep 100`: %s", err)
	}
	syscall.Kill(cmd.Process.Pid, syscall.SIGTERM)
	go ReapOrphans()
	<-testSig
	// Then, do another run to exercise the path for exiting early.
	osGetpid = os.Getpid
	ori := subreaper
	subreaper = func() bool {
		// Call the original subreaper to ensure it gets included in coverage.
		ori()
		return false
	}
	syscall.Kill(syscall.Getpid(), syscall.SIGCHLD)
	go ReapOrphans()
	<-testSig
}
