package helper

import (
	"reflect"
	"strings"
)

func RemoveField(obj interface{}, fields ...string) map[string]interface{} {
	result := make(map[string]interface{})
	rt, rv := reflect.TypeOf(obj), reflect.ValueOf(obj)
	for i := 0; i < rt.NumField(); i++ {
		field := rt.Field(i)
		toRemove := false
		for e := 0; e < len(fields); e++ {
			if strings.EqualFold(field.Name, fields[e]) {
				toRemove = true
				break
			}
		}

		if !toRemove {
			switch rv.Field(i).Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				result[field.Name] = rv.Field(i).Int()
			case reflect.String:
				result[field.Name] = rv.Field(i).String()
			case reflect.Bool:
				result[field.Name] = rv.Field(i).Bool()
			case reflect.Struct:
				result[field.Name] = rv.Field(i)
			}
		}
	}

	return result
}
