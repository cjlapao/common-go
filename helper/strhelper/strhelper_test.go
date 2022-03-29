package strhelper

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToStringArray(t *testing.T) {

	input := "!carlos,test@,abc#123"

	obj, err := ToStringArray(input)

	// Expected result
	assert.Nilf(t, err, "The error should be nil")
	assert.NotNilf(t, obj, "The result not be nil")
	assert.Lenf(t, obj, 3, "There should be 3 strings in the array")
}
