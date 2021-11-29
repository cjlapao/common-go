package security

import (
	"encoding/base64"

	"github.com/cjlapao/common-go/guard"
)

func DecodeBase64String(value string) (string, error) {
	isEmptyKey := guard.EmptyOrNil(value)
	if isEmptyKey != nil {
		return "", isEmptyKey
	}
	decodedBytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}

	return string(decodedBytes), nil
}

func EncodeString(value string) (string, error) {
	isEmptyKey := guard.EmptyOrNil(value)
	if isEmptyKey != nil {
		return "", isEmptyKey
	}

	encodedBytes := base64.StdEncoding.EncodeToString([]byte(value))

	return string(encodedBytes), nil
}
