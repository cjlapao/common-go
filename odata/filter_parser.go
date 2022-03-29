package odata

import "github.com/cjlapao/common-go/parser"

// GlobalFilterTokenizer the global filter tokenizer
var globalFilterTokenizer = filterTokenizer()

// GlobalFilterParser the global filter parser
var globalFilterParser = filterParser()

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
