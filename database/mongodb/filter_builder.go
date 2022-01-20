package mongodb

import (
	"encoding/hex"
	"strings"

	"github.com/cjlapao/common-go/parser"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GlobalFilterTokenizer the global filter tokenizer
var globalFilterTokenizer = filterTokenizer()

// GlobalFilterParser the global filter parser
var globalFilterParser = filterParser()

type FilterParser struct {
	filter string
}

func NewFilterParser(filter string) FilterParser {
	return FilterParser{
		filter: filter,
	}
}

func (filterParser FilterParser) Parse() (interface{}, error) {
	parsedFilter, err := parseFilterString(filterParser.filter)
	if err != nil {
		return nil, err
	}

	result, err := applyFilter(parsedFilter)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// ParseFilterString Converts an input string from the $filter part of the URL into a parse
// tree that can be used by providers to create a response.
func parseFilterString(filter string) (*parser.ParseNode, error) {
	tokens, err := globalFilterTokenizer.Tokenize(filter)
	if err != nil {
		return nil, err
	}

	tree, err := globalFilterParser.Parse(tokens)
	if err != nil {
		return nil, err
	}

	return tree, nil
}

// FilterTokenizer Creates a tokenizer capable of tokenizing filter statements
func filterTokenizer() *parser.Tokenizer {
	tokenizer := parser.Tokenizer{}
	tokenizer.Add("^\\(", parser.FilterTokenOpenParen)
	tokenizer.Add("^\\)", parser.FilterTokenCloseParen)
	tokenizer.Add("^,", parser.FilterTokenComma)
	tokenizer.Add("^(eq|ne|gt|ge|lt|le|and|or) ", parser.FilterTokenLogical)
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

	return &tokenizer
}

// FilterParser creates the definitions for operators and functions
func filterParser() *parser.Parser {
	filterParser := parser.EmptyParser()
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

	return filterParser
}

func (filterParser FilterParser) applyFilter(node *parser.ParseNode) (bson.M, error) {

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

		case "and":
			leftFilter, err := applyFilter(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, _ := applyFilter(node.Children[1]) // Right children
			if err != nil {
				return nil, err
			}
			filter["$and"] = []bson.M{leftFilter, rightFilter}

		case "or":
			leftFilter, err := applyFilter(node.Children[0]) // Left children
			if err != nil {
				return nil, err
			}
			rightFilter, err := applyFilter(node.Children[1]) // Right children
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
