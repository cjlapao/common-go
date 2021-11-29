package helper

import (
	"reflect"
)

func IsNilOrEmpty(value interface{}) bool {
	typeName := reflect.ValueOf(value).Kind().String()

	switch typeName {
	case "string":
		return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
	case "struct":
		return reflect.DeepEqual(value, reflect.Zero(reflect.TypeOf(value)).Interface())
	case "invalid":
		return true
	default:
		return false
	}
}
