package mongodb

type FieldOperation struct {
	Operator MongoOperator
	Value    interface{}
}

type FieldFilterBuilder struct {
	Builder    *FilterBuilder
	FieldName  string
	Operations []FieldOperation
}

func (b *FilterBuilder) Field(fieldName string) *FieldFilterBuilder {
	if len(b.Operations) == 0 {
		b.addSimpleOperation(AND)
	}
	builder := FieldFilterBuilder{
		Builder:   b,
		FieldName: fieldName,
	}

	return &builder
}

func (c *FieldFilterBuilder) Equals(value interface{}) *FieldFilterBuilder {
	c.addFieldOperation(EQUALS, value)
	return c
}

func (c *FieldFilterBuilder) GreaterThan(value interface{}) *FieldFilterBuilder {
	c.addFieldOperation(GREATERTHAN, value)
	return c
}

func (c *FieldFilterBuilder) LowerThan(value interface{}) *FieldFilterBuilder {
	c.addFieldOperation(LOWERTHAN, value)
	return c
}

func (c *FieldFilterBuilder) Build() *FilterBuilder {
	mapper := make(map[string]interface{})
	operationsMapper := make(map[string]interface{})
	for _, operation := range c.Operations {
		operationsMapper[operation.Operator.String()] = operation.Value
	}
	mapper[c.FieldName] = operationsMapper

	return c.Builder
}

func (c *FieldFilterBuilder) addFieldOperation(op MongoOperator, value interface{}) {
	operation := FieldOperation{
		Operator: op,
		Value:    value,
	}
	c.Operations = append(c.Operations, operation)
}

func (c *FieldFilterBuilder) And() *FilterBuilder {
	c.Builder.addSimpleOperation(AND)
	return c.Builder
}
