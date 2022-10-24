package helper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// FileExists Checks if a file/directory exists
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// DirectoryExists Checks if a directory exists
func DirectoryExists(folderPath string) bool {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateDirectory(folderPath string, mode fs.FileMode) bool {
	if err := os.Mkdir(folderPath, mode); err != nil {
		log.Print(fmt.Sprintf("Folder %v could no be created, err: %v", folderPath, err))
		return false
	}

	return true
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
	switch GetOperatingSystem() {
	case WindowsOs:
		return strings.ReplaceAll(path, "/", "\\")
	case LinuxOs:
		return strings.ReplaceAll(path, "\\", "/")
	}

	return ""
}

// JoinPath combines strings into a full path
func JoinPath(items ...string) string {
	var path string
	for _, p := range items {
		if p != "" {
			if len(path) > 0 {
				path += "/"
			}
			path += p
		}
	}
	path = ToOsPath(path)
	// logger.Debug("JoinPath:", path)
	return path
}

func WriteToFile(value string, filePath string) error {
	varBytes := []byte(value)
	f, err := os.Create(filePath)

	defer f.Close()

	if err != nil {
		return err
	}

	f.Write(varBytes)

	f.Sync()

	f.Close()
	return nil
}

func ReadFromFile(filePath string) ([]byte, error) {
	if !FileExists(filePath) {
		return make([]byte, 0), errors.New("File not found")
	}
	return ioutil.ReadFile(filePath)
}

func DeleteFile(filePath string) error {
	if FileExists(filePath) {
		err := os.Remove(filePath)

		if err != nil {
			return err
		}
	}

	return nil
}

func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// _, err = os.Stat(dst)
	// if err != nil && !os.IsNotExist(err) {
	// 	return
	// }
	// if err == nil {
	// 	return fmt.Errorf("destination already exists")
	// }

	// err = os.MkdirAll(dst, si.Mode())
	// if err != nil {
	// 	return
	// }

	if FileExists(src) {
		err = os.MkdirAll(dst, si.Mode())
		if err != nil {
			return
		}
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func DeleteAllFiles(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

func Checksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}

	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}
