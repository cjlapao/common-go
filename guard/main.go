package guard

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/cjlapao/common-go/helper"
)

func FatalEmptyOrNil(value interface{}, name ...string) {
	if helper.IsNilOrEmpty(value) {
		var message string
		if len(name) == 1 {
			message = fmt.Sprintf("Value %v of type %v cannot be nil", name[0], reflect.TypeOf(value))
		} else {
			message = fmt.Sprintf("Value %v cannot be nil", fmt.Sprint(reflect.TypeOf(value)))
		}

		panic(message)
	}
}

func EmptyOrNil(value interface{}, name ...string) error {
	if helper.IsNilOrEmpty(value) {
		var message string
		if len(name) == 1 {
			message = fmt.Sprintf("Value %v of type %v cannot be nil", name[0], reflect.TypeOf(value))
		} else {
			message = fmt.Sprintf("Value %v cannot be nil", fmt.Sprint(reflect.TypeOf(value)))
		}
		return errors.New(message)
	}

	return nil
}

func IsFalse(value bool, name ...string) error {
	if !value {
		var message string
		if len(name) == 1 {
			message = fmt.Sprintf("Value %v cannot be false", name[0])
		} else {
			message = "Value bool cannot be false"
		}

		return errors.New(message)
	}

	return nil
}

func IsNill(value interface{}) bool {
	switch v := value.(type) {
	case string:
		return v == ""
	case bool:
		return false
	case int, int8, int16, int32, int64:
		return false
	case float32, float64:
		return false
	case uint, uintptr, uint8, uint16, uint32, uint64:
		return false
	default:
		if reflect.ValueOf(value).Kind() == reflect.Struct {
			return false
		}
		return value == nil || reflect.ValueOf(value).IsNil()
	}
}
