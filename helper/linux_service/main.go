package linux_service

import (
	"strings"

	"github.com/cjlapao/common-go/commands"
)

type LinuxServiceState int

const (
	LinuxServiceRunning LinuxServiceState = iota
	LinuxServiceStopped
	LinuxServiceErrored
	LinuxServiceUnknown
)

func Status(svcName string) LinuxServiceState {
	output, err := commands.Execute("service", svcName, "status")

	if err != nil {
		return LinuxServiceErrored
	}

	if strings.Contains(output.GetAllOutputs(), "is not running") {
		return LinuxServiceStopped
	}

	if strings.Contains(output.GetAllOutputs(), "is running") {
		return LinuxServiceRunning
	}

	return LinuxServiceUnknown
}

func Start(svcName string) error {
	_, err := commands.Execute("service", svcName, "start")

	if err != nil {
		return err
	}

	return nil
}

func Stop(svcName string) error {
	_, err := commands.Execute("service", svcName, "stop")

	if err != nil {
		return err
	}

	return nil
}

func Restart(svcName string) error {
	_, err := commands.Execute("service", svcName, "restart")

	if err != nil {
		return err
	}

	return nil
}
