package restapiclient

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func DownloadFile(request *http.Request, key string) (*ApiUploadedFile, error) {
	var file multipart.File
	var header *multipart.FileHeader
	var err error

	if err = request.ParseMultipartForm(0); err != nil {
		return nil, err
	}

	if file, header, err = request.FormFile("file"); err != nil {
		return nil, err
	}

	uploadedFile := ApiUploadedFile{
		MultipartFile:       &file,
		MultipartFileHeader: header,
	}

	defer file.Close()

	var fileContentBuffer bytes.Buffer
	io.Copy(&fileContentBuffer, file)

	uploadedFile.File = fileContentBuffer.Bytes()
	uploadedFile.Name = header.Filename
	uploadedFile.Size = header.Size

	return &uploadedFile, nil
}
