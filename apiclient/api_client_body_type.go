package apiclient

import (
	"bytes"
	"encoding/json"
)

type ApiClientBodyType int

const (
	BODY_TYPE_NONE ApiClientBodyType = iota + 1
	BODY_TYPE_JSON
	BODY_TYPE_TEXT
	BODY_TYPE_HTML
	BODY_TYPE_FORM_DATA
	BODY_TYPE_X_WWW_FORM_URLENCODED
	GRAPHQL
)

func (s ApiClientBodyType) String() string {
	return toApiClientBodyTypeString[s]
}

func (s ApiClientBodyType) FromString(key string) ApiClientBodyType {
	return toApiClientBodyTypeID[key]
}

var toApiClientBodyTypeString = map[ApiClientBodyType]string{
	BODY_TYPE_NONE:                  "NONE",
	BODY_TYPE_JSON:                  "JSON",
	BODY_TYPE_TEXT:                  "TEXT",
	BODY_TYPE_HTML:                  "HTML",
	BODY_TYPE_FORM_DATA:             "FORM-DATA",
	BODY_TYPE_X_WWW_FORM_URLENCODED: "X-WWW-FORM-URLENCODED",
	GRAPHQL:                         "GRAPHQL",
}

var toApiClientBodyTypeID = map[string]ApiClientBodyType{
	"NONE":                  BODY_TYPE_NONE,
	"JSON":                  BODY_TYPE_JSON,
	"TEXT":                  BODY_TYPE_TEXT,
	"HTML":                  BODY_TYPE_HTML,
	"FORM-DATA":             BODY_TYPE_FORM_DATA,
	"X-WWW-FORM-URLENCODED": BODY_TYPE_X_WWW_FORM_URLENCODED,
	"GRAPHQL":               GRAPHQL,
}

func (s ApiClientBodyType) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toApiClientBodyTypeString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *ApiClientBodyType) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toApiClientBodyTypeID[j]
	return nil
}

func (s ApiClientBodyType) MarshalYAML() (interface{}, error) {
	return toApiClientBodyTypeString[s], nil
}

func (s *ApiClientBodyType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	j := ""
	err := unmarshal(&j)
	if err != nil {
		return err
	}

	*s = toApiClientBodyTypeID[j]
	return nil
}

func (s *ApiClientBodyType) GetHeader() (key string, value string) {
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
