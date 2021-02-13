package commands

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
)

func Execute(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	var err bytes.Buffer
	cmd.Stderr = &err
	var in bytes.Buffer
	cmd.Stdin = &in

	cmd.Start()
	cmd.Wait()

	errString := err.String()
	if len(errString) > 0 {
		return out.String(), errors.New(errString)
	}

	return out.String(), nil
}

func ExecuteAndWatch(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()
	var out bytes.Buffer
	cmd.Stdout = &out
	var err bytes.Buffer
	cmd.Stderr = &err
	var in bytes.Buffer
	cmd.Stdin = &in

	cmd.Start()

	outScanner := bufio.NewScanner(stdout)
	for outScanner.Scan() {
		m := outScanner.Text()
		fmt.Println(m)
	}

	errorScanner := bufio.NewScanner(stderr)
	for errorScanner.Scan() {
		m := errorScanner.Text()
		fmt.Println(m)
	}

	cmd.Wait()

	errString := err.String()
	if len(errString) > 0 {
		return out.String(), errors.New(errString)
	}

	return out.String(), nil
}
