package apiclient

import (
	"bytes"
	"encoding/json"
)

type ApiClientBodyType int

const (
	NONE ApiClientBodyType = iota + 1
	JSON
	TEXT
	HTML
	FORM_DATA
	X_WWW_FORM_URLENCODED
	GRAPHQL
)

func (s ApiClientBodyType) String() string {
	return toApiClientBodyTypeString[s]
}

func (s ApiClientBodyType) FromString(key string) ApiClientBodyType {
	return toApiClientBodyTypeID[key]
}

var toApiClientBodyTypeString = map[ApiClientBodyType]string{
	NONE:                  "NONE",
	JSON:                  "JSON",
	TEXT:                  "TEXT",
	HTML:                  "HTML",
	FORM_DATA:             "FORM-DATA",
	X_WWW_FORM_URLENCODED: "X-WWW-FORM-URLENCODED",
	GRAPHQL:               "GRAPHQL",
}

var toApiClientBodyTypeID = map[string]ApiClientBodyType{
	"NONE":                  NONE,
	"JSON":                  JSON,
	"TEXT":                  TEXT,
	"HTML":                  HTML,
	"FORM-DATA":             FORM_DATA,
	"X-WWW-FORM-URLENCODED": X_WWW_FORM_URLENCODED,
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

func (s ApiClientBodyType) MarshalYAML() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	buffer.WriteString(toApiClientBodyTypeString[s])
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

func (s *ApiClientBodyType) UnmarshalYAML(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}

	*s = toApiClientBodyTypeID[j]
	return nil
}
