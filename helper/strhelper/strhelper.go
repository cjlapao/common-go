package strhelper

import (
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/guard"
	"github.com/pkg/errors"
)

var ErrEmptyValue = errors.New("field value cannot be empty")
var ErrEmptySeparator = errors.New("separator cannot be empty")
var ErrEmptyArray = errors.New("cannot parse zero length string")

func ToBoolean(value string) bool {
	switch strings.ToLower(value) {
	case "true", "t", "1":
		return true
	case "false", "f", "0":
		return false
	default:
		return false
	}
}

func ToInt(value string) (int, error) {
	result, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil {
		return 0, err
	}
	return result, nil
}

func ToStringArray(value string) ([]string, error) {
	return ToStringArrayF(value, ",")
}

func ToStringArrayF(value string, separator string) ([]string, error) {
	if err := guard.EmptyOrNil(value, "value"); err != nil {
		return nil, ErrEmptyValue
	}
	if err := guard.EmptyOrNil(separator, "separator"); err != nil {
		return nil, ErrEmptySeparator
	}

	// Separating the string into array
	result := strings.Split(value, separator)

	if len(result) == 0 {
		return nil, ErrEmptyArray
	}

	// Triming potential spaces between the
	for idx, resultNoSpace := range result {
		result[idx] = strings.TrimSpace(resultNoSpace)
	}

	return result, nil
}
