package goexec

import (
	"context"
	"os"
	"os/exec"
	"time"
)

const waitDelay = 1000 * time.Hour

func setCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = waitDelay
}

// Command creates a new command with the given context, name, and arguments.
// It sets useful defaults for stdin, stdout, stderr, and signal handling.
func Command(ctx context.Context, name string, arg ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, name, arg...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	setCancel(cmd)
	return cmd
}
