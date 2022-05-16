package apiclient

import (
	"bytes"
	"net/url"
)

type ApiClientBody struct {
	Type     ApiClientBodyType
	Files    []string
	FormData url.Values
	Raw      []byte
}

func NewApiClientBody() *ApiClientBody {
	result := ApiClientBody{
		Type:     JSON,
		Raw:      make([]byte, 0),
		FormData: make(url.Values),
		Files:    make([]string, 0),
	}

	return &result
}

func NewApiClientJsonBody(value []byte) *ApiClientBody {
	result := NewApiClientBody()

	result.Raw = value

	return result
}

func (body *ApiClientBody) Get() *bytes.Buffer {
	switch body.Type {
	case JSON:
		return bytes.NewBuffer(body.Raw)
	case X_WWW_FORM_URLENCODED:
		if len(body.FormData) > 0 {
			body.Raw = []byte(body.FormData.Encode())
			return bytes.NewBuffer(body.Raw)
		} else {
			return nil
		}
	default:
		return bytes.NewBuffer(body.Raw)
	}
}

func (body *ApiClientBody) WithFormValue(key string, value string) *ApiClientBody {
	if body.Type != X_WWW_FORM_URLENCODED {
		body.Type = X_WWW_FORM_URLENCODED
	}

	if !body.FormData.Has(key) {
		body.FormData.Add(key, value)
	} else {
		body.FormData.Set(key, value)
	}

	return body
}
