package restapiclient

import (
	"bytes"
	"encoding/json"
)

type RestApiClientBodyType int

const (
	BODY_TYPE_NONE RestApiClientBodyType = iota + 1
	BODY_TYPE_JSON
	BODY_TYPE_TEXT
	BODY_TYPE_HTML
	BODY_TYPE_FORM_DATA
	BODY_TYPE_X_WWW_FORM_URLENCODED
	GRAPHQL
)

func (s RestApiClientBodyType) String() string {
	return toRestApiClientBodyTypeString[s]
}

func (s RestApiClientBodyType) FromString(key string) RestApiClientBodyType {
	return toRestApiClientBodyTypeID[key]
}

var toRestApiClientBodyTypeString = map[RestApiClientBodyType]string{
	BODY_TYPE_NONE:                  "NONE",
	BODY_TYPE_JSON:                  "JSON",
	BODY_TYPE_TEXT:                  "TEXT",
	BODY_TYPE_HTML:                  "HTML",
	BODY_TYPE_FORM_DATA:             "FORM-DATA",
	BODY_TYPE_X_WWW_FORM_URLENCODED: "X-WWW-FORM-URLENCODED",
	GRAPHQL:                         "GRAPHQL",
}

var toRestApiClientBodyTypeID = map[string]RestApiClientBodyType{
	"NONE":                  BODY_TYPE_NONE,
	"JSON":                  BODY_TYPE_JSON,
	"TEXT":                  BODY_TYPE_TEXT,
	"HTML":                  BODY_TYPE_HTML,
	"FORM-DATA":             BODY_TYPE_FORM_DATA,
	"X-WWW-FORM-URLENCODED": BODY_TYPE_X_WWW_FORM_URLENCODED,
	"GRAPHQL":               GRAPHQL,
}

func (s RestApiClientBodyType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toRestApiClientBodyTypeString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *RestApiClientBodyType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toRestApiClientBodyTypeID[j]
	return nil
}

func (s RestApiClientBodyType) MarshalYAML() (interface{}, error) {
	return toRestApiClientBodyTypeString[s], nil
}

func (s *RestApiClientBodyType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	j := ""
	err := unmarshal(&j)
	if err != nil {
		return err
	}

	*s = toRestApiClientBodyTypeID[j]
	return nil
}

func (s *RestApiClientBodyType) GetHeader() (key string, value string) {
	key = "Content-Type"
	switch *s {
	case BODY_TYPE_JSON:
		value = "application/json;charset=UTF-8"
	case BODY_TYPE_TEXT:
		value = "plain/text"
	case BODY_TYPE_HTML:
		value = "text/html"
	case BODY_TYPE_FORM_DATA:
		value = "multipart/form-data"
	case BODY_TYPE_X_WWW_FORM_URLENCODED:
		value = "application/x-www-form-urlencoded"
	}

	return key, value
}
