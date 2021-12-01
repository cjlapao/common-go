package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RemoveFieldTestStructSubType struct {
	SomeProperty string
}

type RemoveFieldTestStruct struct {
	ID   string
	Name string
	Age  int32
	Sub  RemoveFieldTestStructSubType
}

func TestRemoveFieldShouldNotContainTargetField(t *testing.T) {
	testObj := RemoveFieldTestStruct{
		ID:   "shouldnotbethere",
		Name: "Test",
		Age:  44,
		Sub: RemoveFieldTestStructSubType{
			SomeProperty: "test",
		},
	}

	result := RemoveField(testObj, "ID")

	assert.Equal(t, "test", result["Sub"])
	assert.Empty(t, result["ID"])
}
