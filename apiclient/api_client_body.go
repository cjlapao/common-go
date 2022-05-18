package apiclient

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/url"
	"os"

	"github.com/cjlapao/common-go/helper"
)

type ApiClientBody struct {
	Type    ApiClientBodyType
	Files   map[string]map[string][]byte
	Fields  url.Values
	Raw     []byte
	error   error
	mWriter *multipart.Writer
}

func NewApiClientBody() *ApiClientBody {
	result := ApiClientBody{
		Type:   BODY_TYPE_NONE,
		Raw:    make([]byte, 0),
		Fields: make(url.Values),
		Files:  make(map[string]map[string][]byte, 0),
	}

	return &result
}

func NewApiClientJsonBody(value []byte) *ApiClientBody {
	result := NewApiClientBody()
	result.Type = BODY_TYPE_JSON

	result.Raw = value

	return result
}

func (body *ApiClientBody) IsValid() bool {
	if body.error == nil {
		return true
	} else {
		return false
	}
}

func (body *ApiClientBody) Get() *bytes.Buffer {
	switch body.Type {
	case BODY_TYPE_JSON:
		return bytes.NewBuffer(body.Raw)
	case BODY_TYPE_TEXT:
		return bytes.NewBuffer(body.Raw)
	case BODY_TYPE_HTML:
		return bytes.NewBuffer(body.Raw)
	case BODY_TYPE_X_WWW_FORM_URLENCODED:
		return body.processUrlEncodedData()
	case BODY_TYPE_FORM_DATA:
		return body.processFormData()
	default:
		return bytes.NewBuffer(body.Raw)
	}
}

func (body *ApiClientBody) GetHeader() (key string, value string) {
	key = "Content-Type"
	switch body.Type {
	case BODY_TYPE_JSON:
		value = "application/json;charset=UTF-8"
	case BODY_TYPE_TEXT:
		value = "plain/text"
	case BODY_TYPE_HTML:
		value = "text/html"
	case BODY_TYPE_FORM_DATA:
		value = "multipart/form-data;"
		if body.mWriter != nil {
			value = body.mWriter.FormDataContentType()
		}
	case BODY_TYPE_X_WWW_FORM_URLENCODED:
		value = "application/x-www-form-urlencoded"
	}

	return key, value
}

func (body *ApiClientBody) UrlEncoded() *ApiClientBody {
	if body.Type != BODY_TYPE_X_WWW_FORM_URLENCODED {
		body.Type = BODY_TYPE_X_WWW_FORM_URLENCODED
	}

	return body
}

func (body *ApiClientBody) FormData() *ApiClientBody {
	if body.Type != BODY_TYPE_FORM_DATA {
		body.Type = BODY_TYPE_FORM_DATA
	}

	return body
}

func (body *ApiClientBody) Json(content []byte) *ApiClientBody {
	if body.Type != BODY_TYPE_JSON {
		body.Type = BODY_TYPE_JSON
	}

	body.Raw = content
	return body
}

func (body *ApiClientBody) Text(content []byte) *ApiClientBody {
	if body.Type != BODY_TYPE_TEXT {
		body.Type = BODY_TYPE_TEXT
	}

	body.Raw = content
	return body
}

func (body *ApiClientBody) Html(content []byte) *ApiClientBody {
	if body.Type != BODY_TYPE_HTML {
		body.Type = BODY_TYPE_HTML
	}

	body.Raw = content
	return body
}

func (body *ApiClientBody) WithFile(key string, filepath string) *ApiClientBody {
	if helper.FileExists(filepath) {
		var fileContent []byte

		fileContent, err := helper.ReadFromFile(filepath)
		if err != nil {
			body.error = err
			return body
		}

		stat, err := os.Stat(filepath)
		if err != nil {
			body.error = err
			return body
		}

		body.Files[key] = make(map[string][]byte)
		body.Files[key][stat.Name()] = fileContent
	}

	return body
}

func (body *ApiClientBody) WithField(key string, value string) *ApiClientBody {
	if !body.Fields.Has(key) {
		body.Fields.Add(key, value)
	} else {
		body.Fields.Set(key, value)
	}

	return body
}

func (body *ApiClientBody) processUrlEncodedData() *bytes.Buffer {
	if len(body.Fields) == 0 {
		return nil
	}

	body.Raw = []byte(body.Fields.Encode())
	return bytes.NewBuffer(body.Raw)
}

func (body *ApiClientBody) processFormData() *bytes.Buffer {
	if len(body.Fields) == 0 && len(body.Files) == 0 {
		return nil
	}

	var b bytes.Buffer

	body.mWriter = multipart.NewWriter(&b)
	for key, value := range body.Fields {
		if fw, err := body.mWriter.CreateFormField(key); err == nil {
			io.Copy(fw, bytes.NewBuffer([]byte(value[0])))
		}
	}
	for key, files := range body.Files {
		for file, fileContent := range files {
			if fw, err := body.mWriter.CreateFormFile(key, file); err == nil {
				io.Copy(fw, bytes.NewBuffer(fileContent))
			}
		}
	}

	body.mWriter.Close()

	return &b
}
