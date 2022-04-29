package apiclient

import (
	"bytes"
	"encoding/json"
)

type ApiClientMethod int

const (
	GET ApiClientMethod = iota + 1
	POST
	DELETE
	PUT
	PATCH
	OPTIONS
	HEAD
	CONNECT
	TRACE
)

func (s ApiClientMethod) String() string {
	return toApiClientMethodString[s]
}

func (s ApiClientMethod) FromString(key string) ApiClientMethod {
	return toApiClientMethodID[key]
}

var toApiClientMethodString = map[ApiClientMethod]string{
	GET:     "GET",
	POST:    "POST",
	DELETE:  "DELETE",
	PUT:     "PUT",
	PATCH:   "PATCH",
	OPTIONS: "OPTIONS",
	HEAD:    "HEAD",
	CONNECT: "CONNECT",
	TRACE:   "TRACE",
}

var toApiClientMethodID = map[string]ApiClientMethod{
	"GET":     GET,
	"POST":    POST,
	"DELETE":  DELETE,
	"PUT":     PUT,
	"PATCH":   PATCH,
	"OPTIONS": OPTIONS,
	"HEAD":    HEAD,
	"CONNECT": CONNECT,
	"TRACE":   TRACE,
}

func (s ApiClientMethod) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toApiClientMethodString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *ApiClientMethod) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toApiClientMethodID[j]
	return nil
}

func (s ApiClientMethod) MarshalYAML() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toApiClientMethodString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *ApiClientMethod) UnmarshalYAML(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toApiClientMethodID[j]
	return nil
}
