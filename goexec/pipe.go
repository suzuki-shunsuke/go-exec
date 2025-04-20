package goexec

import (
	"errors"
	"io"
	"os/exec"
)

// Pipe executes a series of commands in a pipeline, where the output of each command is passed as input to the next command.
// Pipe returns the last command executed, the index of the command that failed, and an error if any occurred.
func Pipe(cmds ...*exec.Cmd) (*exec.Cmd, int, error) {
	s := len(cmds)
	if s < 2 { //nolint:mnd
		return nil, 0, errors.New("at least two commands are required")
	}
	entries := make([]*entry, s)
	for i, cmd := range cmds {
		entries[i] = &entry{
			cmd: cmd,
		}
	}
	s1 := s - 1
	for i := range s1 {
		reader, writer := io.Pipe()
		entries[i].writer = writer
		setStdout(entries[i].cmd, writer)
		entries[i+1].reader = reader
		entries[i+1].cmd.Stdin = reader
	}

	for i, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return cmd, i, err //nolint:wrapcheck
		}
	}
	for i, entry := range entries {
		if err := entry.Wait(); err != nil {
			return entry.cmd, i, err
		}
	}
	return nil, 0, nil
}

type entry struct {
	cmd    *exec.Cmd
	reader io.Closer
	writer io.Closer
}

func (e *entry) Wait() error {
	defer func() {
		if e.reader != nil {
			e.reader.Close()
		}
		if e.writer != nil {
			e.writer.Close()
		}
	}()
	if err := e.cmd.Wait(); err != nil {
		return err //nolint:wrapcheck
	}
	return nil
}

func setStdout(cmd *exec.Cmd, writer io.Writer) {
	if cmd.Stdout == nil {
		cmd.Stdout = writer
	} else {
		cmd.Stdout = io.MultiWriter(writer, cmd.Stdout)
	}
}
