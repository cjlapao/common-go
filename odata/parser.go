package parser

import (
	"errors"
	"net/url"
	"strings"

	"github.com/cjlapao/common-go/helper/strhelper"
	"github.com/cjlapao/common-go/odata/filterparser"
)

// OData keywords
const (
	Select      = "$select"
	Top         = "$top"
	Skip        = "$skip"
	Count       = "$count"
	OrderBy     = "$orderby"
	InlineCount = "$inlinecount"
	Filter      = "$filter"
)

// ParseURLValues parses url values in odata format into a map of interfaces for the DB adapters to translate
//nolint :gocyclo
func ParseURLValues(query url.Values) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	var parseErrors []string

	result[Count] = false
	result[InlineCount] = "none"

	if isCountAndInlineCountSet(query) {
		parseErrors = append(parseErrors, "$count and $inlinecount cannot be set in the same odata query")
	}

	for queryParam, queryValues := range query {
		var parseResult interface{}
		var err error

		if len(queryValues) > 1 {
			parseErrors = append(parseErrors, "Duplicate keyword '"+queryParam+"' found in odata query")
			continue
		}
		value := query.Get(queryParam)
		if value == "" && queryParam != Count {
			parseErrors = append(parseErrors, "No value was set for keyword '"+queryParam+"'")
			continue
		}

		switch queryParam {
		case Select:
			parseResult, err = strhelper.ToStringArray(value)
		case Top:
			parseResult, err = strhelper.ToInt(value)
		case Skip:
			parseResult, err = strhelper.ToInt(value)
		case Count:
			parseResult = true
		case OrderBy:
			parseResult, err = parseOrderArray(value)
		case InlineCount:
			if !isValidInlineCountValue(value) {
				parseErrors = append(parseErrors, "Inline count value needs to be allpages or none")
			}
			parseResult = value
		case Filter:
			parseResult, err = filterparser.ParseFilterString(value)
		default:
			parseErrors = append(parseErrors, "Keyword '"+queryParam+"' is not valid")
		}

		if err != nil {
			parseErrors = append(parseErrors, err.Error())
		}
		result[queryParam] = parseResult
	}
	if len(parseErrors) > 0 {
		return nil, errors.New(strings.Join(parseErrors[:], ";"))
	}
	return result, nil
}

func isValidInlineCountValue(value string) bool {
	valueNoSpace := strings.TrimSpace(value)
	if valueNoSpace != "allpages" && valueNoSpace != "none" {
		return false
	}
	return true
}

func isCountAndInlineCountSet(query url.Values) bool {

	_, countFound := query[Count]
	_, inlineFound := query[InlineCount]

	if countFound && inlineFound {
		return true
	}

	return false
}
