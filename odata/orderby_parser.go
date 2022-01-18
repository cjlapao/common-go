package parser

import (
	"errors"
	"strings"

	"github.com/cjlapao/common-go/helper/strhelper"
	"github.com/cjlapao/common-go/validators"
)

// OrderItem holds order key information
type OrderItem struct {
	Field string
	Order string
}

const (
	// Ascendent Order
	Ascendent = "asc"
	// Descendent Order
	Descendent = "desc"
)

var ErrToManyOrderByStatements = errors.New("cannot have more than 2 items in orderby query")
var ErrInvalidOrderBy = errors.New("second value in orderby needs to be asc or desc")

func parseOrderArray(value string) ([]OrderItem, error) {
	parsedArray, err := strhelper.ToStringArray(value)
	if err != nil {
		return nil, err
	}

	// Validate values for special characters
	valid := validators.New("~!@#$%^&*()_+-")
	for _, val := range parsedArray {
		if valid.ValidateField(val) || val == "" {
			return nil, errors.New("Cannot support field " + val)
		}
	}

	result := make([]OrderItem, len(parsedArray))

	for i, v := range parsedArray {
		timmedString := strings.TrimSpace(v)
		compressedSpaces := strings.Join(strings.Fields(timmedString), " ")
		s := strings.Split(compressedSpaces, " ")

		if len(s) > 2 {
			return nil, ErrToManyOrderByStatements
		}

		if len(s) > 1 {
			if s[1] != Ascendent &&
				s[1] != Descendent {
				return nil, ErrInvalidOrderBy
			}
			result[i] = OrderItem{s[0], s[1]}
			continue
		}
		result[i] = OrderItem{compressedSpaces, "asc"}
	}
	return result, nil
}
