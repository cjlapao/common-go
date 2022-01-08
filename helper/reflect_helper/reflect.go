package reflect_helper

import (
	"reflect"
	"strings"
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

func GetFieldTag(field reflect.StructField, key string) string {
	fieldTag := string(field.Tag)
	tags := strings.Split(fieldTag, " ")
	for _, tag := range tags {
		tagDetails := strings.Split(tag, ":")
		if len(tagDetails) == 2 {
			tagKey := tagDetails[0]
			if strings.EqualFold(tagKey, key) {
				result := ""
				tagValue := tagDetails[1]
				tagValueDetails := strings.Split(tagValue, ",")
				if len(tagDetails) == 1 {
					result = tagValue
				}
				if len(tagDetails) == 2 {
					result = tagValueDetails[0]
				}
				result = strings.ReplaceAll(result, "\"", "")
				return result
			}
		}
	}

	return ""
}
