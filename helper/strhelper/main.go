package strhelper

import "strings"

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
