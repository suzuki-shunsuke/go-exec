package goexec

import (
	"bytes"
	"io"
	"os/exec"
)

// Runner executes commands and captures their output.
// This is used to mock command execution in tests.
type Runner struct{}

// New creates a new Executor.
func New() *Runner {
	return &Runner{}
}

// Result holds the result of a command execution.
type Result struct {
	ExitCode       int
	Stdout         *bytes.Buffer
	Stderr         *bytes.Buffer
	CombinedOutput *bytes.Buffer
}

// Run executes the command and sets the result.
// If you don't need result, you can pass nil.
func (e *Runner) Run(cmd *exec.Cmd, result *Result) error {
	if result == nil {
		result = &Result{}
	}
	if result.Stdout != nil {
		cmd.Stdout = io.MultiWriter(cmd.Stdout, result.Stdout)
	}
	if result.Stderr != nil {
		cmd.Stderr = io.MultiWriter(cmd.Stderr, result.Stderr)
	}
	if result.CombinedOutput != nil {
		cmd.Stdout = io.MultiWriter(cmd.Stdout, result.CombinedOutput)
		cmd.Stderr = io.MultiWriter(cmd.Stderr, result.CombinedOutput)
	}
	err := cmd.Run()
	result.ExitCode = cmd.ProcessState.ExitCode()
	return err
}
