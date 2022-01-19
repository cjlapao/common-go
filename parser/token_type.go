package parser

// Token constants
const (
	FilterTokenOpenParen int = iota
	FilterTokenCloseParen
	FilterTokenWhitespace
	FilterTokenComma
	FilterTokenLogical
	FilterTokenFunc
	FilterTokenFloat
	FilterTokenInteger
	FilterTokenString
	FilterTokenDate
	FilterTokenTime
	FilterTokenDateTime
	FilterTokenBoolean
	FilterTokenLiteral
)
