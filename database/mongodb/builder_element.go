package mongodb

import "errors"

var ErrNoElements = errors.New("no elements to process")

type elementBuilderOperation int

const (
	SetOperation elementBuilderOperation = iota
	SetOnInsertOperation
	UnsetOperation
	FilterOperation
)

type builderElement struct {
	operation       elementBuilderOperation
	key             string
	filterOperation filterOperation
	value           interface{}
}

type BuilderOptions int

const (
	UpsertBuildOption BuilderOptions = 1
)
