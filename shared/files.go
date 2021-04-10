package shared

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

type Startup struct {
}

// FileExists Checks if a file/directory exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func Delete(path string) error {
	if FileExists(path) {
		err := os.Remove(path)

		if err != nil {
			return err
		}
	}

	return nil
}

// GetExecutionPath gets executable path
func GetExecutionPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return ToOsPath(dir)
}

// ToOsPath Converts a path into the native os path
func ToOsPath(path string) string {
	os := runtime.GOOS
	switch strings.ToLower(os) {
	case "linux":
		return strings.ReplaceAll(path, "\\", "/")
	case "windows":
		return strings.ReplaceAll(path, "/", "\\")
	}
	return path
}

// JoinPath combines strings into a full path
func JoinPath(items ...string) string {
	var path string
	for _, p := range items {
		if len(path) > 0 {
			path += "/"
		}
		path += p
	}
	path = ToOsPath(path)
	return path
}
