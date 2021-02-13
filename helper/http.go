package helper

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

// DownloadFile Downloads a file from a url
func DownloadFile(url string, filepath string) error {
	resp, err := http.Get(url)

	CheckError(err)

	if resp.StatusCode != http.StatusOK {
		return errors.New("Error downloading file from " + url + " with code " + fmt.Sprint(resp.StatusCode))
	}

	defer resp.Body.Close()

	out, err := os.Create(filepath)

	CheckError(err)

	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	return err
}
