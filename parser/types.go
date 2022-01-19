package parser

// Operator constants
const (
	OpAssociationLeft int = iota
	OpAssociationRight
)

// Operator operator structure
type Operator struct {
	Token string
	// Whether the operator is left/right/or not associative
	Association int
	// The number of operands this operator operates on
	Operands int
	// Rank of precedence
	Precedence int
}

// Function function structure
type Function struct {
	Token string
	// The number of parameters this function accepts
	Params int
}

// ParseNode parseNode structure
type ParseNode struct {
	Token    *Token
	Parent   *ParseNode
	Children []*ParseNode
}
