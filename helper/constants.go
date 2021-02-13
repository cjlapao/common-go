package helper

import (
	"runtime"
	"strings"
)

// Helpers Constants
const (
	AlphaNumeric = "1234567890abcdefghijklmnopqrstuvwxyz"
)

// OperatingSystem enum
type OperatingSystem int

// Defines the operating system Enum
const (
	WindowsOs OperatingSystem = iota
	LinuxOs
	UnknownOs
)

func (o OperatingSystem) String() string {
	return [...]string{"Windows", "Linux", "Unknown"}[o]
}

// GetOperatingSystem returns the operating system
func GetOperatingSystem() OperatingSystem {
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		return LinuxOs
	case "windows":
		return WindowsOs
	}
	return UnknownOs
}
