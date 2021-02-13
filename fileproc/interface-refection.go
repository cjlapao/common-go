package fileproc

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetArray Gets an array from an interface address
func GetArray(i interface{}, value string) ([]interface{}, error) {
	prop, err := GetPropertyValue(i, value)
	if err != nil {
		return nil, err
	}
	xType := fmt.Sprintf("%T", prop)
	if xType == "[]interface {}" {
		slice := interfaceSlice(prop)
		if slice == nil {
			return nil, errors.New("There was an error converting the interface to an array")
		}
		return slice, nil
	}

	return nil, errors.New("The property is not of type array")
}

// GetBoolean Gets an boolean value from an interface address
func GetBoolean(i interface{}, value string) (bool, error) {
	propValue, err := GetPropertyValue(i, value)
	if err != nil {
		return false, err
	}

	xType := fmt.Sprintf("%T", propValue)
	if xType == "bool" {
		return propValue.(bool), nil
	}

	return false, errors.New("The property is not of type boolean")
}

// GetString Gets an string value from an interface address
func GetString(i interface{}, value string) (string, error) {
	propValue, err := GetPropertyValue(i, value)
	if err != nil {
		return "", err
	}

	xType := fmt.Sprintf("%T", propValue)
	if xType == "string" {
		return propValue.(string), nil
	}

	return "", errors.New("The property is not of type string")
}

// GetInt Gets an int value from an interface address
func GetInt(i interface{}, value string) (int, error) {
	propValue, err := GetPropertyValue(i, value)
	if err != nil {
		return -1, err
	}

	xType := fmt.Sprintf("%T", propValue)
	if xType == "int" {
		return propValue.(int), nil
	}

	return -1, errors.New("The property is not of type int")
}

// GetFloat32 Gets an float32 value from an interface address
func GetFloat32(i interface{}, value string) (float32, error) {
	propValue, err := GetPropertyValue(i, value)
	if err != nil {
		return -1, err
	}

	xType := fmt.Sprintf("%T", propValue)
	if xType == "float32" {
		return propValue.(float32), nil
	}

	return -1, errors.New("The property is not of type float32")
}

// GetFloat64 Gets an float64 value from an interface address
func GetFloat64(i interface{}, value string) (float64, error) {
	propValue, err := GetPropertyValue(i, value)
	if err != nil {
		return -1, err
	}

	xType := fmt.Sprintf("%T", propValue)
	if xType == "float64" {
		return propValue.(float64), nil
	}

	return -1, errors.New("The property is not of type float64")
}

// GetPropertyValue Gets a property generic value
func GetPropertyValue(i interface{}, value string) (interface{}, error) {
	if value == "" {
		return nil, errors.New("The value cannot be empty")
	}

	properties := strings.Split(value, ".")

	if len(properties) > 0 {
		lastProperty, err := getGenericProperty(i, ".")
		if err != nil {
			return nil, errors.New("Root property was not found in interface")
		}
		for i, property := range properties {
			if i == len(properties)-1 {
				propertyValue, err := getGenericProperty(lastProperty, property)
				if err != nil {
					return nil, errors.New("Property was not found in interface")
				}

				logger.Debug("Found value %v for property %v", property, value)
				xType := fmt.Sprintf("%T", propertyValue)

				switch xType {
				case "[]interface {}":
					return propertyValue, nil
				case "bool":
					return propertyValue.(bool), nil
				case "int":
					return propertyValue.(int), nil
				case "float32":
					return propertyValue.(float32), nil
				case "float64":
					return propertyValue.(float64), nil
				}

				propertyStringValue := propertyValue.(string)
				if intValue, err := strconv.Atoi(propertyStringValue); err == nil {
					return intValue, nil
				}
				if booValue, err := strconv.ParseBool(propertyStringValue); err == nil {
					return booValue, nil
				}

				return propertyValue, nil
			}

			arrayStartPos := strings.Index(property, "[")
			if arrayStartPos > -1 {
				propertyName := property[:arrayStartPos]
				arrayEndPos := strings.Index(property, "]")
				arrayValue, err := getGenericProperty(lastProperty, propertyName)
				if err != nil {
					return nil, errors.New("Property was not found in interface")
				}
				xType := fmt.Sprintf("%T", arrayValue)
				if xType == "[]interface {}" {
					arrIndex := property[arrayStartPos+1 : arrayEndPos]
					index, err := strconv.Atoi(arrIndex)
					if err != nil {
						lastProperty = arrayValue.([]interface{})[0]
					} else {
						lastProperty = arrayValue.([]interface{})[index]
					}
				} else {
					break
				}
			} else {
				lastProperty, _ = getGenericProperty(lastProperty, property)
			}
		}
	}
	return (nil), errors.New("Value does not contain any properties")
}

func getGenericProperty(i interface{}, value string) (interface{}, error) {
	m, err := json.Marshal(&i)
	if err != nil {
		return (nil), err
	}

	var x map[string]interface{}

	err = json.Unmarshal(m, &x)
	if err != nil {
		return (nil), err
	}

	if value == "" || value == "." {
		return x, nil
	}

	prop := x[value]

	if prop == nil {
		return nil, errors.New("Property " + value + " was not found")
	}

	return prop, nil
}

func interfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	// Keep the distinction between nil and empty slice input
	if s.IsNil() {
		return nil
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}
