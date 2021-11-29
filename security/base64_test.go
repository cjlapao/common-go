package security_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/cjlapao/common-go/security"
	"github.com/stretchr/testify/assert"
)

var encodedSomeString = "c29tZVN0cmluZw=="
var encodedSomeWrongString = "xxc29tZVN0cmluZf=="

func TestDecodeBase64String_ReturnCorrectValue(t *testing.T) {
	tests := []struct {
		encoded, decoded string
		shouldFail       bool
		err              error
	}{
		{encodedSomeString, "someString", false, nil},
		{encodedSomeWrongString, "", true, errors.New("illegal base64 data at input byte 16")},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("DecodeBase64String Return Correct Value %v", tt.decoded)
		t.Run(testName, func(t *testing.T) {
			decoded, err := security.DecodeBase64String(tt.encoded)

			if tt.shouldFail {
				assert.EqualErrorf(t, err, tt.err.Error(), "Decoding should generate an error")
				assert.Equal(t, "", decoded)
			} else {
				assert.Nilf(t, err, "Error should be nil")
				assert.Equal(t, tt.decoded, decoded)
			}
		})
	}
}

func TestDecodeBase64String_WithEmptyValue_ReturnError(t *testing.T) {
	decoded, err := security.DecodeBase64String("")

	assert.EqualErrorf(t, err, "Value string cannot be nil", "Decoding should generate an error")
	assert.Equal(t, "", decoded)
}
