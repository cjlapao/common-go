package restapiclient

import (
	"bytes"
	"encoding/json"
)

type RestApiClientMethod int

const (
	API_METHOD_GET RestApiClientMethod = iota + 1
	API_METHOD_POST
	API_METHOD_DELETE
	API_METHOD_PUT
	API_METHOD_PATCH
	API_METHOD_OPTIONS
	API_METHOD_HEAD
	API_METHOD_CONNECT
	API_METHOD_TRACE
)

func (s RestApiClientMethod) String() string {
	return toRestApiClientMethodString[s]
}

func (s RestApiClientMethod) FromString(key string) RestApiClientMethod {
	return toRestApiClientMethodID[key]
}

var toRestApiClientMethodString = map[RestApiClientMethod]string{
	API_METHOD_GET:     "GET",
	API_METHOD_POST:    "POST",
	API_METHOD_DELETE:  "DELETE",
	API_METHOD_PUT:     "PUT",
	API_METHOD_PATCH:   "PATCH",
	API_METHOD_OPTIONS: "OPTIONS",
	API_METHOD_HEAD:    "HEAD",
	API_METHOD_CONNECT: "CONNECT",
	API_METHOD_TRACE:   "TRACE",
}

var toRestApiClientMethodID = map[string]RestApiClientMethod{
	"GET":     API_METHOD_GET,
	"POST":    API_METHOD_POST,
	"DELETE":  API_METHOD_DELETE,
	"PUT":     API_METHOD_PUT,
	"PATCH":   API_METHOD_PATCH,
	"OPTIONS": API_METHOD_OPTIONS,
	"HEAD":    API_METHOD_HEAD,
	"CONNECT": API_METHOD_CONNECT,
	"TRACE":   API_METHOD_TRACE,
}

func (s RestApiClientMethod) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toRestApiClientMethodString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *RestApiClientMethod) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toRestApiClientMethodID[j]
	return nil
}

func (s RestApiClientMethod) MarshalYAML() (interface{}, error) {
	return toRestApiClientMethodString[s], nil
}

func (s *RestApiClientMethod) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var j string
	err := unmarshal(&j)
	if err != nil {
		return err
	}

	*s = toRestApiClientMethodID[j]
	return nil
}
