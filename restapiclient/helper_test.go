package restapiclient

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cjlapao/common-go/helper"
	"github.com/stretchr/testify/assert"
)

func TestHelper_DownloadFile_WithCorrectResponse(t *testing.T) {
	testFileName := "test.Unit"
	testFileContent := "someFileContent"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := DownloadFile(r, "file")
		fileSize := len([]byte(testFileContent))
		assert.NotNilf(t, file, "file should not be nil")
		assert.Nilf(t, err, "error should be nil")
		assert.Equalf(t, testFileName, file.Name, "incorrect file name. want %v, found %v", testFileName, file.Name)
		assert.Equalf(t, int64(fileSize), file.Size, "incorrect file size. want %v, found %v", fileSize, file.Size)
		assert.Equalf(t, testFileContent, string(file.File), "incorrect file content. want %v, found %v", testFileContent, string(file.File))
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	api.SendRequest(RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody().FormData().WithFile("file", testFileName),
	})

	helper.DeleteFile(testFileName)
}

func TestHelper_DownloadFile_WithNoUploadedFile_NilIsReturned(t *testing.T) {
	testFileName := "test.Unit"
	testFileContent := "someFileContent"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := DownloadFile(r, "file")
		assert.Nilf(t, file, "file should not be nil")
		assert.NotNilf(t, err, "error should be nil")
		assert.Equal(t, "request Content-Type isn't multipart/form-data", err.Error())
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	api.SendRequest(RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody(),
	})

	helper.DeleteFile(testFileName)
}

func Test_Helper_DownloadFile_WithNoFileKey_NilIsReturned(t *testing.T) {
	testFileName := "test.Unit"
	testFileContent := "someFileContent"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		file, err := DownloadFile(r, "file")
		assert.Nilf(t, file, "file should not be nil")
		assert.NotNilf(t, err, "error should be nil")
		assert.Equal(t, "http: no such file", err.Error())
	}))

	defer server.Close()

	api := DefaultRestApiClient{
		Client: server.Client(),
	}

	if helper.FileExists(testFileName) {
		helper.DeleteFile(testFileName)
	}

	helper.WriteToFile(testFileContent, testFileName)

	api.SendRequest(RestApiClientRequest{
		Method: API_METHOD_POST,
		Body:   NewRestApiClientBody().FormData().WithFile("other_file", testFileName),
	})

	helper.DeleteFile(testFileName)
}
