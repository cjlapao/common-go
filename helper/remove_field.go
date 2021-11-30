package helper

import (
	"reflect"
	"strings"
)

func RemoveField(obj *interface{}, fields ...string) map[string]interface{} {
	result := make(map[string]interface{})
	rt, rv := reflect.TypeOf(*obj), reflect.ValueOf(*obj)
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
			result[field.Name] = rv. .String()
		}
	}
}
