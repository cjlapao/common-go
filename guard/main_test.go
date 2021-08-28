package guard

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	TestString string
}

func TestFatalEmptyOrNil(t *testing.T) {
	emptyStruc := TestStruct{}
	var nilStruc TestStruct
	var nilInterface interface{}
	emptyString := ""
	nonEmptyString := "foo"
	nonEmptyStruct := TestStruct{
		TestString: "bar",
	}

	assert.PanicsWithValuef(t, "Value guard.TestStruct cannot be nil", func() { FatalEmptyOrNil(emptyStruc) }, "Value should issue a panic with a specific message")
	assert.PanicsWithValuef(t, "Value guard.TestStruct cannot be nil", func() { FatalEmptyOrNil(nilStruc) }, "Value should issue a panic with a specific message")
	assert.PanicsWithValuef(t, "Value <nil> cannot be nil", func() { FatalEmptyOrNil(nilInterface) }, "Value should issue a panic with a specific message")
	assert.PanicsWithValuef(t, "Value string cannot be nil", func() { FatalEmptyOrNil(emptyString) }, "Value should issue a panic with a specific message")

	assert.PanicsWithValuef(t, "Value emptyStruc of type guard.TestStruct cannot be nil", func() { FatalEmptyOrNil(emptyStruc, "emptyStruc") }, "Value should issue a panic with a specific message")
	assert.PanicsWithValuef(t, "Value nilStruc of type guard.TestStruct cannot be nil", func() { FatalEmptyOrNil(nilStruc, "nilStruc") }, "Value should issue a panic with a specific message")
	assert.PanicsWithValuef(t, "Value nilInterface of type <nil> cannot be nil", func() { FatalEmptyOrNil(nilInterface, "nilInterface") }, "Value should issue a panic with a specific message")
	assert.PanicsWithValuef(t, "Value emptyString of type string cannot be nil", func() { FatalEmptyOrNil(emptyString, "emptyString") }, "Value should issue a panic with a specific message")

	assert.NotPanics(t, func() { FatalEmptyOrNil(nonEmptyString) })
	assert.NotPanics(t, func() { FatalEmptyOrNil(nonEmptyStruct) })
}

func TestEmptyOrNil(t *testing.T) {
	// Arrange
	emptyStruc := TestStruct{}
	var nilStruc TestStruct
	var nilInterface interface{}

	emptyString := ""
	nonEmptyString := "foo"
	nonEmptyStruct := TestStruct{
		TestString: "bar",
	}

	// Act + Assert
	assert.EqualErrorf(t, EmptyOrNil(emptyStruc), "Value guard.TestStruct cannot be nil", "Empty Struct should issue an error")
	assert.EqualErrorf(t, EmptyOrNil(nilStruc), "Value guard.TestStruct cannot be nil", "nil Struct should issue an error")
	assert.EqualErrorf(t, EmptyOrNil(nilInterface), "Value <nil> cannot be nil", "nil interface should issue an error")
	assert.EqualErrorf(t, EmptyOrNil(emptyString), "Value string cannot be nil", "Empty string should issue an error")

	assert.EqualErrorf(t, EmptyOrNil(emptyStruc, "emptyStruc"), "Value emptyStruc of type guard.TestStruct cannot be nil", "Empty Struct should issue an error with variable name")
	assert.EqualErrorf(t, EmptyOrNil(nilStruc, "nilStruc"), "Value nilStruc of type guard.TestStruct cannot be nil", "nil Struct should issue an error with variable name")
	assert.EqualErrorf(t, EmptyOrNil(nilInterface, "nilInterface"), "Value nilInterface of type <nil> cannot be nil", "nil interface should issue an error with variable name")
	assert.EqualErrorf(t, EmptyOrNil(emptyString, "emptyString"), "Value emptyString of type string cannot be nil", "Empty string should issue an error with variable name")

	assert.NoErrorf(t, EmptyOrNil(nonEmptyString), "Non Empty string should not issue an error")
	assert.NoErrorf(t, EmptyOrNil(nonEmptyStruct), "Non Empty struct should not issue an error")
}

func TestIsFalse(t *testing.T) {
	// Arrange
	trueValue := true
	falseValue := false

	// Act + Assert
	assert.EqualErrorf(t, IsFalse(falseValue), "Value bool cannot be false", "Empty Struct should issue an error")
	assert.EqualErrorf(t, IsFalse(falseValue, "falseValue"), "Value falseValue cannot be false", "Empty Struct should issue an error")
	assert.NoErrorf(t, IsFalse(trueValue), "Non Empty struct should not issue an error")
}
