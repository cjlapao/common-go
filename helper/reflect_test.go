package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNilOrEmpty(t *testing.T) {
	// Arrange
	emptyString := ""
	nonEmptyString := "foo"
	zeroVal := 0
	falseVal := false
	int64Val := int64(0)
	floatVal := float64(0)
	emptyStructValue := TestStructure{}
	var nilStructValue TestStructure
	var nilInterfaceValue interface{}
	nonEmptyStruct := TestStructure{
		TestString: "bar",
	}

	// Act + Assert
	assert.True(t, IsNilOrEmpty(nilInterfaceValue))
	assert.True(t, IsNilOrEmpty(emptyString))
	assert.False(t, IsNilOrEmpty(nonEmptyString))
	assert.False(t, IsNilOrEmpty(zeroVal))
	assert.False(t, IsNilOrEmpty(falseVal))
	assert.False(t, IsNilOrEmpty(int64Val))
	assert.False(t, IsNilOrEmpty(floatVal))
	assert.True(t, IsNilOrEmpty(emptyStructValue))
	assert.True(t, IsNilOrEmpty(nilStructValue))
	assert.False(t, IsNilOrEmpty(nonEmptyStruct))
}
