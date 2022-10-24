package commands

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
)

func Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var stdOut, stdIn bytes.Buffer

	cmd.Stdout = io.MultiWriter(os.Stdout, &stdOut)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stdOut)
	cmd.Stdin = &stdIn

	if err := cmd.Run(); err != nil {
		return stdOut.String(), err
	}

	return stdOut.String(), nil
}

func ExecuteAndWatch(command string, args ...string) (string, error) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	var stdOut, stdIn bytes.Buffer

	cmd.Stdout = io.MultiWriter(os.Stdout, &stdOut)
	cmd.Stderr = io.MultiWriter(os.Stderr, &stdOut)
	cmd.Stdin = &stdIn

	if err := cmd.Start(); err != nil {
		return stdOut.String(), err
	}

	go func() {
		<-ctx.Done()
	}()

	if err := cmd.Wait(); err != nil {
		return stdOut.String(), err
	}

	return stdOut.String(), nil
}
