package helper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNilOrEmpty(t *testing.T) {
	var p *string
	v := reflect.ValueOf(p)

	assert.Equal(t, true, v.IsValid())
	assert.True(t, IsNilOrEmpty(v))

	assert.Equal(t, uintptr(0), v.Pointer())

	v = v.Elem()
	assert.Equal(t, false, v.IsValid())
	assert.True(t, IsNilOrEmpty(v))
}
