package helper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type RemoveFieldStruct struct {
	ID          string
	Name        string
	IsOld       bool
	Age         int
	AndSomeMore interface{}
	OneMap      map[string]interface{}
	Sub         RemoveFieldSub
}

type RemoveFieldSub struct {
	Hello string
}

func TestRemoveField(t *testing.T) {
	// Arrange
	test := RemoveFieldStruct{
		ID:          "ID",
		Name:        "SomeName",
		IsOld:       true,
		Age:         20,
		AndSomeMore: "yep",
		OneMap:      make(map[string]interface{}),
		Sub: RemoveFieldSub{
			Hello: "world",
		},
	}
	test.OneMap["testing"] = "something"

	result := RemoveField(test, "ID")

	assert.Equal(t, result["Name"], "SomeName")
}
