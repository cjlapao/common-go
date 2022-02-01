package mongodb

type elementBuilderOperation int

const (
	SetOperation elementBuilderOperation = iota
	SetOnInsertOperation
	UnsetOperation
	FilterOperation
)

type BuilderElement struct {
	operation       elementBuilderOperation
	key             string
	filterOperation filterOperation
	value           interface{}
}
