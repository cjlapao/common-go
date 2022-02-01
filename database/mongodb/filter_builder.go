package mongodb

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/cjlapao/common-go/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type filter struct {
	field     string
	operation filterOperation
	value     interface{}
}

type filterOperation string

const (
	GreaterThan        filterOperation = "gt"
	GreaterOrEqualThan filterOperation = "ge"
	LowerThan          filterOperation = "lt"
	LowerOrEqualThan   filterOperation = "le"
	Equal              filterOperation = "eq"
	NotEqual           filterOperation = "ne"
	Contains           filterOperation = "contains"
	StartsWith         filterOperation = "startswith"
	EndsWith           filterOperation = "endswith"
	RegEx              filterOperation = "regex"
)

type FilterParser struct {
	globalFilterTokenizer *parser.Tokenizer
	globalFilterParser    *parser.Parser
	filter                string
}

// NewFilterParser Creates a nem Filter Parser from a odata type of query
func NewFilterParser(filter string) *FilterParser {
	result := FilterParser{
		filter: filter,
	}

	result.globalFilterParser = result.filterParser()
	result.globalFilterTokenizer = result.filterTokenizer()

	return &result
}

// Parse Creates a MongoDB compatible filter from a odata type of query
func (filterParser *FilterParser) Parse() (interface{}, error) {
	parsedFilter, err := filterParser.parseFilterString(filterParser.filter)
	if err != nil {
		return nil, err
	}

	result, err := ApplyFilter(parsedFilter)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// parseFilterString Converts an input string from a base odata type of query into a parse
// tree that can be used by providers to create a compatible mongodb filter query.
func (fp FilterParser) parseFilterString(filter string) (*parser.ParseNode, error) {
	tokens, err := fp.globalFilterTokenizer.Tokenize(filter)
	if err != nil {
		return nil, err
	}

	tree, err := fp.globalFilterParser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// filterTokenizer Creates a tokenizer capable of tokenizing filter statements
func (fp *FilterParser) filterTokenizer() *parser.Tokenizer {
	tokenizer := parser.Tokenizer{}
	tokenizer.Add("^\\(", parser.FilterTokenOpenParen)
	tokenizer.Add("^\\)", parser.FilterTokenCloseParen)
	tokenizer.Add("^,", parser.FilterTokenComma)
	tokenizer.Add("^(eq|ne|gt|ge|lt|le|and|or|regex) ", parser.FilterTokenLogical)
	tokenizer.Add("^(contains|endswith|startswith)", parser.FilterTokenFunc)
	tokenizer.Add("^-?[0-9]+\\.[0-9]+", parser.FilterTokenFloat)
	tokenizer.Add("^-?[0-9]+", parser.FilterTokenInteger)
	tokenizer.Add("^(?i:true|false)", parser.FilterTokenBoolean)
	tokenizer.Add("^'(''|[^'])*'", parser.FilterTokenString)
	tokenizer.Add("^-?[0-9]{4,4}-[0-9]{2,2}-[0-9]{2,2}", parser.FilterTokenDate)
	tokenizer.Add("^[0-9]{2,2}:[0-9]{2,2}(:[0-9]{2,2}(.[0-9]+)?)?", parser.FilterTokenTime)
	tokenizer.Add("^[0-9]{4,4}-[0-9]{2,2}-[0-9]{2,2}T[0-9]{2,2}:[0-9]{2,2}(:[0-9]{2,2}(.[0-9]+)?)?(Z|[+-][0-9]{2,2}:[0-9]{2,2})", parser.FilterTokenDateTime)
	tokenizer.Add("^[a-zA-Z][a-zA-Z0-9_.]*", parser.FilterTokenLiteral)
	tokenizer.Add("^_id", parser.FilterTokenLiteral)
	tokenizer.Ignore("^ ", parser.FilterTokenWhitespace)

	fp.globalFilterTokenizer = &tokenizer
	return fp.globalFilterTokenizer
}

// filterParser creates the definitions for operators and functions
func (fp *FilterParser) filterParser() *parser.Parser {
	filterParser := parser.EmptyParser()
	filterParser.DefineOperator("regex", 2, parser.OpAssociationLeft, 4)
	filterParser.DefineOperator("gt", 2, parser.OpAssociationLeft, 4)
	filterParser.DefineOperator("ge", 2, parser.OpAssociationLeft, 4)
	filterParser.DefineOperator("lt", 2, parser.OpAssociationLeft, 4)
	filterParser.DefineOperator("le", 2, parser.OpAssociationLeft, 4)
	filterParser.DefineOperator("eq", 2, parser.OpAssociationLeft, 3)
	filterParser.DefineOperator("ne", 2, parser.OpAssociationLeft, 3)
	filterParser.DefineOperator("and", 2, parser.OpAssociationLeft, 2)
	filterParser.DefineOperator("or", 2, parser.OpAssociationLeft, 1)
	filterParser.DefineFunction("contains", 2)
	filterParser.DefineFunction("endswith", 2)
	filterParser.DefineFunction("startswith", 2)

	fp.globalFilterParser = filterParser
	return fp.globalFilterParser
}

func ApplyFilter(node *parser.ParseNode) (bson.M, error) {

	filter := make(bson.M)

	if _, ok := node.Token.Value.(string); ok {
		switch node.Token.Value {

		case "eq":
			// Escape single quotes in the case of strings
			if _, valueOk := node.Children[1].Token.Value.(string); valueOk {
				node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			}
			value := bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, keyOk := node.Children[0].Token.Value.(string); !keyOk {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "ne":
			// Escape single quotes in the case of strings
			if _, valueOk := node.Children[1].Token.Value.(string); valueOk {
				node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			}
			value := bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, keyOk := node.Children[0].Token.Value.(string); !keyOk {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "gt":
			var keyString string
			if keyString, ok = node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}

			var value bson.M
			if keyString == "_id" {
				var idString string
				if _, ok := node.Children[1].Token.Value.(string); ok {
					idString = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
				}
				decodedString, err := hex.DecodeString(idString)
				if err != nil || len(decodedString) != 12 {
					return nil, ErrInvalidInput
				}
				objectId := primitive.NewObjectID()
				objectId.UnmarshalText(decodedString)
				value = bson.M{"$" + node.Token.Value.(string): objectId}
			} else {
				value = bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			}
			filter[keyString] = value

		case "ge":
			value := bson.M{"$gte": node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "lt":
			value := bson.M{"$" + node.Token.Value.(string): node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "le":
			value := bson.M{"$lte": node.Children[1].Token.Value}
			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "regex":
			// Escape single quotes in the case of strings
			if _, valueOk := node.Children[1].Token.Value.(string); valueOk {
				node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)
			}
			// value := bson.M{"$regex": node.Children[1].Token.Value, "$options": "i"}
			value := primitive.Regex{
				Pattern: node.Children[1].Token.Value.(string),
				Options: "gi",
			}

			if _, ok := node.Children[0].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "and":
			leftFilter, err := ApplyFilter(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, _ := ApplyFilter(node.Children[1]) // Right children
			if err != nil {
				return nil, err
			}
			filter["$and"] = []bson.M{leftFilter, rightFilter}

		case "or":
			leftFilter, err := ApplyFilter(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, err := ApplyFilter(node.Children[1]) // Right children
			if err != nil {
				return nil, err
			}
			filter["$or"] = []bson.M{leftFilter, rightFilter}

		//Functions
		case "startswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)

			value := primitive.Regex{
				Pattern: "^" + node.Children[1].Token.Value.(string),
				Options: "gi",
			}

			filter[node.Children[0].Token.Value.(string)] = value

		case "endswith":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)

			value := primitive.Regex{
				Pattern: node.Children[1].Token.Value.(string) + "$",
				Options: "gi",
			}
			filter[node.Children[0].Token.Value.(string)] = value

		case "contains":
			if _, ok := node.Children[1].Token.Value.(string); !ok {
				return nil, ErrInvalidInput
			}
			node.Children[1].Token.Value = strings.Replace(node.Children[1].Token.Value.(string), "'", "", -1)

			value := primitive.Regex{
				Pattern: node.Children[1].Token.Value.(string),
				Options: "gi",
			}
			filter[node.Children[0].Token.Value.(string)] = value

		}
	}
	return filter, nil
}

func getOperationString(field string, operation filterOperation, value interface{}) string {
	switch operation {
	case "gt":
		return fmt.Sprintf("%v gt '%v'", field, value)
	case "ge":
		return fmt.Sprintf("%v ge '%v'", field, value)
	case "lt":
		return fmt.Sprintf("%v lt '%v'", field, value)
	case "le":
		return fmt.Sprintf("%v le '%v'", field, value)
	case "eq":
		return fmt.Sprintf("%v eq '%v'", field, value)
	case "ne":
		return fmt.Sprintf("%v ne '%v'", field, value)
	case "regex":
		return fmt.Sprintf("%v regex %v", field, value)
	case "contains":
		return fmt.Sprintf("contains(%v, '%v')", field, value)
	case "endswith":
		return fmt.Sprintf("endswith(%v, '%v')", field, value)
	case "startswith":
		return fmt.Sprintf("startswith(%v, '%v')", field, value)
	default:
		return fmt.Sprintf("%v eq '%v'", field, value)
	}
}
