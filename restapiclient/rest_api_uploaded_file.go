package restapiclient

import "mime/multipart"

type ApiUploadedFile struct {
	File                []byte
	Name                string
	Size                int64
	MultipartFile       *multipart.File
	MultipartFileHeader *multipart.FileHeader
}
