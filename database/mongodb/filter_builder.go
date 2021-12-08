package mongodb

import (
	"encoding/json"
	"strings"
)

type FilterBuilder struct {
	Operations []FilterOperation
}

type FilterOperation struct {
	LeftOperator  MongoOperator
	Elements      []FilterElement
	TestElements  []FieldFilterBuilder
	RightOperator MongoOperator
}

func NewFilterBuilder() *FilterBuilder {
	return &FilterBuilder{
		Operations: make([]FilterOperation, 0),
	}
}

func (b *FilterBuilder) And() *FilterBuilder {
	b.addSimpleOperation(AND)
	return b
}

func (b *FilterBuilder) OldAnd(fieldName string, fieldValue string) *FilterBuilder {
	b.addOperation(AND, fieldName, fieldValue)
	return b
}

func (b *FilterBuilder) Or(fieldName string, fieldValue string) *FilterBuilder {
	b.addOperation(OR, fieldName, fieldValue)
	return b
}

func (b *FilterBuilder) Nor(fieldName string, fieldValue string) *FilterBuilder {
	b.addOperation(NOR, fieldName, fieldValue)
	return b
}

func (b *FilterBuilder) Build() interface{} {
	result := make(map[string]interface{})
	for _, operation := range b.Operations {
		elementMaps := make([]map[string]interface{}, 0)
		for _, element := range operation.Elements {
			elementMaps = append(elementMaps, element.Encode())
		}
		result[operation.LeftOperator.String()] = elementMaps
	}

	j, _ := json.Marshal(result)
	test := string(j)
	println(test)
	return result
}

func (b *FilterBuilder) addOperation(op MongoOperator, fieldName string, fieldValue interface{}) {
	element := FilterElement{
		Key:   fieldName,
		Value: fieldValue,
	}

	if len(b.Operations) == 0 {
		operation := FilterOperation{
			LeftOperator:  op,
			Elements:      make([]FilterElement, 0),
			RightOperator: NONE,
		}
		operation.Elements = append(operation.Elements, element)
		b.Operations = append(b.Operations, operation)
		return
	}

	exists, idx := b.operationExists(op)
	if exists {
		operation := &b.Operations[idx]
		if idx > 0 {
			previousOperation := &b.Operations[idx-1]
			previousOperation.RightOperator = op
		}
		operation.Elements = append(operation.Elements, element)
	} else {
		previousOperation := &b.Operations[len(b.Operations)-1]
		previousOperation.RightOperator = op
		operation := FilterOperation{
			LeftOperator:  op,
			Elements:      make([]FilterElement, 0),
			RightOperator: NONE,
		}
		operation.Elements = append(operation.Elements, element)
		b.Operations = append(b.Operations, operation)
	}
}

func (b *FilterBuilder) addSimpleOperation(op MongoOperator) {
	if len(b.Operations) == 0 {
		operation := FilterOperation{
			LeftOperator:  op,
			Elements:      make([]FilterElement, 0),
			RightOperator: NONE,
		}
		b.Operations = append(b.Operations, operation)
		return
	}

	exists, idx := b.operationExists(op)
	if exists {
		if idx > 0 {
			previousOperation := &b.Operations[idx-1]
			previousOperation.RightOperator = op
		}
	} else {
		previousOperation := &b.Operations[len(b.Operations)-1]
		previousOperation.RightOperator = op
		operation := FilterOperation{
			LeftOperator:  op,
			Elements:      make([]FilterElement, 0),
			RightOperator: NONE,
		}
		b.Operations = append(b.Operations, operation)
	}
}

func (b *FilterBuilder) operationExists(op MongoOperator) (bool, int) {
	for idx, operation := range b.Operations {
		if operation.LeftOperator == op {
			return true, idx
		}
	}

	return false, -1
}

func (b *FilterBuilder) fieldExistsInOperation(op MongoOperator, fieldName string) (bool, int) {
	for _, operation := range b.Operations {
		if operation.LeftOperator == op {
			for idx, field := range operation.TestElements {
				if strings.EqualFold(field.FieldName, fieldName) {
					return true, idx
				}
			}
		}
	}

	return false, -1
}

func (b *FilterBuilder) getCurrentOperation() *FilterOperation {
	if len(b.Operations) == 0 {
		operation := FilterOperation{
			LeftOperator:  NONE,
			Elements:      make([]FilterElement, 0),
			RightOperator: NONE,
		}
		b.Operations = append(b.Operations, operation)
		return &operation
	} else {
		return &b.Operations[len(b.Operations)-1]
	}
}
