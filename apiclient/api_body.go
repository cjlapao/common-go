package apiclient

import "net/url"

type ApiClientBody struct {
	Type     ApiClientBodyType
	Files    []string
	FormData url.Values
	Raw      []byte
}
