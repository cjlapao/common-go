package automapper

// Based in the Peter Strøiman automapper with some additions

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/cjlapao/common-go/helper/reflect_helper"
	"github.com/cjlapao/common-go/helper/strhelper"
)

type MapOptions int64

const (
	Loose                  MapOptions = 1
	RequestForm            MapOptions = 2
	RequestFormWithJsonTag MapOptions = 3
)

// Map fills out the fields in dest with values from source. All fields in the
// destination object must exist in the source object.
//
// Object hierarchies with nested structs and slices are supported, as long as
// type types of nested structs/slices follow the same rules, i.e. all fields
// in destination structs must be found on the source struct.
//
// Embedded/anonymous structs are supported
//
// Values that are not exported/not public will not be mapped.
//
// It is a design decision to panic when a field cannot be mapped in the
// destination to ensure that a renamed field in either the source or
// destination does not result in subtle silent bug.
func Map(source, dest interface{}, options ...MapOptions) error {
	var err error
	var destType = reflect.TypeOf(dest)
	if destType.Kind() != reflect.Ptr {
		return errors.New("dest must be a pointer type")
	}

	var sourceVal = reflect.ValueOf(source)
	var destVal = reflect.ValueOf(dest).Elem()
	if len(options) == 0 {
		err = mapValues(sourceVal, destVal, false)
	}
	for _, option := range options {
		switch option {
		case Loose:
			err = mapValues(sourceVal, destVal, true)
		case RequestForm:
			err = mapRequestForm(source, dest, "")
			return err
		case RequestFormWithJsonTag:
			err = mapRequestForm(source, dest, "json")
			return err
		}
	}

	return err
}

func mapRequestForm(source interface{}, dest interface{}, tag string) error {
	switch r := source.(type) {
	case *http.Request:
		var form url.Values

		err := r.ParseForm()
		if err == nil && len(r.Form) > 0 {
			form = r.Form
		}

		if len(form) == 0 {
			err = r.ParseMultipartForm(0)
			if err == nil && len(r.Form) > 0 {
				form = r.Form
			}
		}
		if len(form) == 0 {
			return errors.New("no data found to map")
		}

		if tag != "" {
			mapFormValues(form, dest, true, tag)
		} else {
			mapFormValues(form, dest, false, tag)
		}
	default:
		return errors.New("source must be a pointer http request")
	}

	return nil
}

func mapFormValues(sourceVal url.Values, dest interface{}, useTag bool, tagName string) {
	var destVal = reflect.ValueOf(dest).Elem()
	destType := destVal.Type()
	for i := 0; i < destType.NumField(); i++ {
		fieldName := destType.Field(i).Name
		fieldTag := reflect_helper.GetFieldTag(destType.Field(i), "json")
		var sourceFieldVal string
		if useTag {
			if tagName == "" {
				tagName = "json"
			}
			sourceFieldVal = sourceVal.Get(fieldTag)
		} else {
			sourceFieldVal = sourceVal.Get(fieldName)
			// trying the value in lowercase
			if sourceFieldVal == "" {
				sourceFieldVal = sourceVal.Get(strings.ToLower(fieldName))
			}
		}
		switch destType.Field(i).Type.Kind() {
		case reflect.Bool:
			boolValue := strhelper.ToBoolean(sourceFieldVal)
			destVal.Field(i).SetBool(boolValue)
		case reflect.String:
			destVal.Field(i).SetString(sourceFieldVal)
		case reflect.Float32, reflect.Float64:
			if s, err := strconv.ParseFloat(sourceFieldVal, 64); err == nil {
				destVal.Field(i).SetFloat(s)
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if s, err := strconv.ParseInt(sourceFieldVal, 10, 64); err == nil {
				destVal.Field(i).SetInt(s)
			}
		}
	}
}

func mapValues(sourceVal, destVal reflect.Value, loose bool) error {
	var err error
	destType := destVal.Type()
	if destType.Kind() == reflect.Struct {
		if sourceVal.Type().Kind() == reflect.Ptr {
			if sourceVal.IsNil() {
				// If source is nil, it maps to an empty struct
				sourceVal = reflect.New(sourceVal.Type().Elem())
			}
			sourceVal = sourceVal.Elem()
		}
		for i := 0; i < destVal.NumField(); i++ {
			err = mapField(sourceVal, destVal, i, loose)
			if err != nil {
				return err
			}
		}
	} else if destType == sourceVal.Type() {
		destVal.Set(sourceVal)
	} else if destType.Kind() == reflect.Ptr {
		if valueIsNil(sourceVal) {
			return nil
		}
		val := reflect.New(destType.Elem())
		err = mapValues(sourceVal, val.Elem(), loose)
		if err != nil {
			return err
		}
		destVal.Set(val)
	} else if destType.Kind() == reflect.Slice {
		err := mapSlice(sourceVal, destVal, loose)
		if err != nil {
			return err
		}
	} else {
		return errors.New("currently not supported")
	}
	return err
}

func mapSlice(sourceVal, destVal reflect.Value, loose bool) error {
	var err error
	destType := destVal.Type()
	length := sourceVal.Len()
	target := reflect.MakeSlice(destType, length, length)
	for j := 0; j < length; j++ {
		val := reflect.New(destType.Elem()).Elem()
		err = mapValues(sourceVal.Index(j), val, loose)
		if err != nil {
			return err
		}
		target.Index(j).Set(val)
	}

	if length == 0 {
		err = verifyArrayTypesAreCompatible(sourceVal, destVal, loose)
	}
	destVal.Set(target)
	return err
}

func verifyArrayTypesAreCompatible(sourceVal, destVal reflect.Value, loose bool) error {
	dummyDest := reflect.New(reflect.PtrTo(destVal.Type()))
	dummySource := reflect.MakeSlice(sourceVal.Type(), 1, 1)
	err := mapValues(dummySource, dummyDest.Elem(), loose)
	return err
}

func mapField(source, destVal reflect.Value, i int, loose bool) error {
	var err error
	destType := destVal.Type()
	fieldName := destType.Field(i).Name
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("error mapping field: %s. DestType: %v. SourceType: %v. Error: %v", fieldName, destType, source.Type(), r)
			panic(err.Error())
		}
	}()

	destField := destVal.Field(i)
	if destType.Field(i).Anonymous {
		err = mapValues(source, destField, loose)
	} else {
		if valueIsContainedInNilEmbeddedType(source, fieldName) {
			return nil
		}
		sourceField := source.FieldByName(fieldName)
		if (sourceField == reflect.Value{}) {
			if loose {
				return nil
			}
			if destField.Kind() == reflect.Struct {
				err := mapValues(source, destField, loose)
				if err != nil {
					return err
				}
				return nil
			} else {
				for i := 0; i < source.NumField(); i++ {
					if source.Field(i).Kind() != reflect.Struct {
						continue
					}
					if sourceField = source.Field(i).FieldByName(fieldName); (sourceField != reflect.Value{}) {
						break
					}
				}
			}
		}
		err = mapValues(sourceField, destField, loose)
	}

	return err
}

func valueIsNil(value reflect.Value) bool {
	return value.Type().Kind() == reflect.Ptr && value.IsNil()
}

func valueIsContainedInNilEmbeddedType(source reflect.Value, fieldName string) bool {
	structField, _ := source.Type().FieldByName(fieldName)
	ix := structField.Index
	if len(structField.Index) > 1 {
		parentField := source.FieldByIndex(ix[:len(ix)-1])
		if valueIsNil(parentField) {
			return true
		}
	}
	return false
}
