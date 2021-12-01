package executionctx

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/cjlapao/common-go/helper"
	"github.com/stretchr/testify/assert"
)

func TestGetConfigurationProvider_ReturnCorrect(t *testing.T) {
	// Arrange + Act
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	key1 := config.Get("foo")
	initializedConfig := GetConfigService()

	// Assert
	assert.False(t, helper.IsNilOrEmpty(config))
	assert.True(t, helper.IsNilOrEmpty(key1))
	assert.False(t, helper.IsNilOrEmpty(initializedConfig))
}

func TestGetConfigurationProvider_ResturnsSameKeyAfterReinitialization(t *testing.T) {
	// Arrange + Act
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	key1 := config.Get("foo")
	config.UpsertKey("foo", "bar")
	updatedKey1 := config.Get("foo")
	initializedConfig := GetConfigService()
	key2 := initializedConfig.Get("foo")

	// Assert
	assert.False(t, helper.IsNilOrEmpty(config))
	assert.True(t, helper.IsNilOrEmpty(key1))
	assert.False(t, helper.IsNilOrEmpty(initializedConfig))
	assert.Equal(t, "bar", updatedKey1)
	assert.Equal(t, "bar", key2)
}

func TestConfigurationProvider_IfUpsertEmptyKeyErrorIsReturned(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()

	//Act
	key1 := config.UpsertKey("", "bar")

	//Assert
	assert.IsTypef(t, errors.New("someError"), key1, "Config Provider should return an error interface")
}

func TestConfigurationProvider_IfUpsertEmptyValueErrorIsReturned(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	var emptyInterface interface{}
	emptyStruct := helper.TestStructure{}

	//Act
	key1 := config.UpsertKey("foo", "")
	key2 := config.UpsertKey("bar", nil)
	key3 := config.UpsertKey("interface", emptyInterface)
	key4 := config.UpsertKey("struc", emptyStruct)

	key1Value := config.Get("foo")
	key2Value := config.Get("bar")
	key3Value := config.Get("interface")
	key4Value := config.Get("struc")

	//Assert
	assert.IsTypef(t, errors.New("someError"), key1, "Config Provider should return an error interface")
	assert.IsTypef(t, errors.New("someError"), key2, "Config Provider should return an error interface")
	assert.IsTypef(t, errors.New("someError"), key3, "Config Provider should return an error interface")
	assert.IsTypef(t, errors.New("someError"), key4, "Config Provider should return an error interface")

	assert.Truef(t, helper.IsNilOrEmpty(key1Value), "Config Provider should return empty value for key1")
	assert.Truef(t, helper.IsNilOrEmpty(key2Value), "Config Provider should return empty value for key2")
	assert.Truef(t, helper.IsNilOrEmpty(key3Value), "Config Provider should return empty value for key3")
	assert.Truef(t, helper.IsNilOrEmpty(key4Value), "Config Provider should return empty value for key4")
}

func TestConfigurationProvider_IfUpsertKeyValueIsStoredInVault(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	testStruct := helper.TestStructure{
		TestString: "bar",
	}

	//Act
	key1 := config.UpsertKey("foo", "bar")
	key2 := config.UpsertKey("bar", 2)
	key3 := config.UpsertKey("interface", int64(4))
	key4 := config.UpsertKey("struc", testStruct)

	key1Value := config.Get("foo")
	key2Value := config.Get("bar")
	key3Value := config.Get("interface")
	key4Value := config.Get("struc")

	//Assert
	assert.IsTypef(t, nil, key1, "Config Provider should return an error interface")
	assert.IsTypef(t, nil, key2, "Config Provider should return an error interface")
	assert.IsTypef(t, nil, key3, "Config Provider should return an error interface")
	assert.IsTypef(t, nil, key4, "Config Provider should return an error interface")

	assert.Falsef(t, helper.IsNilOrEmpty(key1Value), "Config Provider should return value for key")
	assert.Falsef(t, helper.IsNilOrEmpty(key2Value), "Config Provider should return value for key")
	assert.Falsef(t, helper.IsNilOrEmpty(key3Value), "Config Provider should return value for key")
	assert.Falsef(t, helper.IsNilOrEmpty(key4Value), "Config Provider should return value for key")

	assert.Equalf(t, "bar", key1Value, "Config Provider should return \"bar\" value for key")
	assert.Equalf(t, 2, key2Value, "Config Provider should return \"2\" value for key")
	assert.Equalf(t, int64(4), key3Value, "Config Provider should return \"4\" value for key")
	assert.Equalf(t, testStruct, key4Value, "Config Provider should return \"struct\" value for key")
}

func TestConfigurationProvider_GetKeyShouldReturnEnvValues(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})
	config := GetConfigService()
	os.Setenv("foo", "bar")

	//Act
	key1Value := config.Get("foo")

	//Assert
	assert.Falsef(t, helper.IsNilOrEmpty(key1Value), "Config Provider should return value for key")

	assert.Equalf(t, "bar", key1Value, "Config Provider should return \"bar\" value for key")
	os.Setenv("foo", "")
}

func TestConfigurationProvider_GetKeyShouldReturnEnvValuesInsteadOfVaultValues(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	os.Setenv("foo", "bar")
	config.UpsertKey("foo", "noPriorityBar")

	//Act
	key1Value := config.Get("foo")

	//Assert
	assert.Falsef(t, helper.IsNilOrEmpty(key1Value), "Config Provider should return value for key")

	assert.Equalf(t, "bar", key1Value, "Config Provider should return \"bar\" value for key")
	os.Setenv("foo", "")
}

func TestConfigurationProvider_UpsertEmptyKeysShouldReturnNil(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	inserts := make(map[string]interface{})

	//Act
	keys := config.UpsertKeys(inserts)

	//Assert
	assert.Nilf(t, keys, "Config Provider should return nill")
}

func TestConfigurationProvider_UpsertEmptyKeysShouldReturnError(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	var inserts map[string]interface{}

	//Act
	keys := config.UpsertKeys(inserts)

	//Assert
	assert.Truef(t, len(keys) == 1, "Config Provider should return nill")
	assert.EqualErrorf(t, keys[0], "array is nil", "Error Message does not match with expected")
}

func TestConfigurationProvider_UpsertKeysShouldAddArrayIntoVault(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	inserts := make(map[string]interface{})
	inserts["foo"] = "bar"
	inserts["on"] = "theMoney"

	//Act
	keys := config.UpsertKeys(inserts)

	//Assert
	assert.Truef(t, len(keys) == 0, "Config Provider should not return an error interface")

	assert.Falsef(t, helper.IsNilOrEmpty(config.Get("foo")), "Config Provider should return value for foo")
	assert.Falsef(t, helper.IsNilOrEmpty(config.Get("on")), "Config Provider should return value for on")

	assert.Equalf(t, "bar", config.Get("foo"), "Config Provider should return \"bar\" value for key foo")
	assert.Equalf(t, "theMoney", config.Get("on"), "Config Provider should return \"bar\" value for key on")
}

func TestConfigurationProvider_UpsertKeysWithErrorsShouldReturnArray(t *testing.T) {
	//Arrange
	// resetting the internal values
	configurationService = nil
	vault = make(map[string]interface{})

	config := GetConfigService()
	inserts := make(map[string]interface{})
	inserts["foo"] = "bar"
	inserts["on"] = "theMoney"
	inserts[""] = "Error"
	inserts["error"] = ""

	//Act
	keys := config.UpsertKeys(inserts)

	//Assert
	assert.Truef(t, len(keys) == 2, "Config Provider should not return an error interface")

	assert.Falsef(t, helper.IsNilOrEmpty(config.Get("foo")), "Config Provider should return value for foo")
	assert.Falsef(t, helper.IsNilOrEmpty(config.Get("on")), "Config Provider should return value for on")

	assert.Equalf(t, "bar", config.Get("foo"), "Config Provider should return \"bar\" value for key foo")
	assert.Equalf(t, "theMoney", config.Get("on"), "Config Provider should return \"bar\" value for key on")

	for _, error := range keys {
		assert.IsTypef(t, errors.New("someError"), error, "Config Provider should return an error interface")

	}
}

func TestNewConfigServiceResetsVault(t *testing.T) {
	//Arrange
	config := NewConfigService()
	config.UpsertKey("foo", "bar")
	keyValue := config.Get("foo")
	config = NewConfigService()
	newKeyValue := config.Get("foo")

	//Assert
	assert.NotNilf(t, config, "Config service should not be nil")
	assert.Equalf(t, "bar", keyValue, "Key value should be \"bar\"")
	assert.True(t, helper.IsNilOrEmpty(newKeyValue), "Key \"foo\" should have been nil or empty after reset but found %v", newKeyValue)
}

func TestClearEmptiesVault(t *testing.T) {
	//Arrange
	config := NewConfigService()
	config.UpsertKey("foo", "bar")

	//Act
	keyValue := config.Get("foo")
	config.Clear()
	keyValueAfterClear := config.Get("foo")

	//assert
	assert.NotNilf(t, config, "Config service should not be nil")
	assert.Equalf(t, "bar", keyValue, "Key value should be \"bar\"")
	assert.Nilf(t, keyValueAfterClear, "Key value should have been clear")
}

func TestGetString_ReturnCorrectValue(t *testing.T) {
	//Arrange
	var tests = []struct {
		varName       string
		value         interface{}
		expectedValue string
	}{
		{"someString", "foo", "foo"},
		{"someInt", 1, "1"},
		{"someBool", true, "true"},
		{"someStruct", helper.TestStructure{TestString: "someMoreString"}, "{someMoreString false 0}"},
	}

	// Act on table
	for _, tt := range tests {
		testName := fmt.Sprintf("Getting String Return Correct Value -> %v,%v", tt.value, tt.varName)
		t.Run(testName, func(t *testing.T) {
			config := NewConfigService()
			config.UpsertKey(tt.varName, tt.value)

			// act
			keyValue := config.GetString(tt.varName)

			assert.Equal(t, tt.expectedValue, keyValue)
		})
	}
}

func TestGetString_WithNotFoundValue_ReturnEmptyString(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetString("notFound")

	assert.Equal(t, "", keyValue)
}

func TestGetString_WithEmptyKey_ReturnEmptyString(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetString("")

	assert.Equal(t, "", keyValue)
}

func TestGetInt_ReturnCorrectValue(t *testing.T) {
	//Arrange
	var tests = []struct {
		varName       string
		value         interface{}
		expectedValue int
	}{
		{"someString", "foo", 0},
		{"someInt", 1, 1},
		{"someFloat", 1.3, 0},
		{"someBool", true, 0},
		{"someStruct", helper.TestStructure{TestString: "someMoreString"}, 0},
	}

	// Act on table
	for _, tt := range tests {
		testName := fmt.Sprintf("Getting Int Return Correct Value -> %v,%v", tt.value, tt.varName)
		t.Run(testName, func(t *testing.T) {
			config := NewConfigService()
			config.UpsertKey(tt.varName, tt.value)

			// act
			keyValue := config.GetInt(tt.varName)

			assert.Equal(t, tt.expectedValue, keyValue)
		})
	}
}

func TestGetInt_WithNotFoundValue_ReturnEmptyString(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetInt("notFound")

	assert.Equal(t, 0, keyValue)
}

func TestGetInt_WithEmptyKey_ReturnEmptyString(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetInt("")

	assert.Equal(t, 0, keyValue)
}

func TestGetBool_ReturnCorrectValue(t *testing.T) {
	//Arrange
	var tests = []struct {
		varName       string
		value         interface{}
		expectedValue bool
	}{
		{"someString", "foo", false},
		{"someBoolFalseString", "f", false},
		{"someBoolFalseString1", "F", false},
		{"someBoolFalseString2", "false", false},
		{"someBoolTrueString", "t", true},
		{"someBooltrueString1", "T", true},
		{"someBoolTrueString2", "true", true},
		{"someInt", 1, true},
		{"someFalseInt", 0, false},
		{"someBool", true, true},
		{"someStruct", helper.TestStructure{TestString: "someMoreString"}, false},
	}

	// Act on table
	for _, tt := range tests {
		testName := fmt.Sprintf("Getting Bool Return Correct Value -> %v,%v", tt.value, tt.varName)
		t.Run(testName, func(t *testing.T) {
			config := NewConfigService()
			config.UpsertKey(tt.varName, tt.value)

			// act
			keyValue := config.GetBool(tt.varName)

			assert.Equal(t, tt.expectedValue, keyValue)
		})
	}
}

func TestGetBool_WithNotFoundValue_ReturnFalse(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetBool("notFound")

	assert.Equal(t, false, keyValue)
}

func TestGetBool_WithEmptyKey_ReturnFalse(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetBool("")

	assert.Equal(t, false, keyValue)
}

func TestGetFloat_ReturnCorrectValue(t *testing.T) {
	//Arrange
	var tests = []struct {
		varName       string
		value         interface{}
		expectedValue float64
	}{
		{"someString", "foo", 0},
		{"someInt", 1, 1},
		{"someFloat", 1.3, 1.3},
		{"someBool", true, 0},
		{"someStruct", helper.TestStructure{TestString: "someMoreString"}, 0},
	}

	// Act on table
	for _, tt := range tests {
		testName := fmt.Sprintf("Getting Int Return Correct Value -> %v,%v", tt.value, tt.varName)
		t.Run(testName, func(t *testing.T) {
			config := NewConfigService()
			config.UpsertKey(tt.varName, tt.value)

			// act
			keyValue := config.GetFloat(tt.varName)

			assert.Equal(t, tt.expectedValue, keyValue)
		})
	}
}

func TestGeFloat_WithNotFoundValue_ReturnZero(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetFloat("notFound")

	assert.Equal(t, float64(0), keyValue)
}

func TestGetFloat_WithEmptyKey_ReturnZero(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetFloat("")

	assert.Equal(t, float64(0), keyValue)
}

func TestGetBase64_ReturnCorrectValue(t *testing.T) {
	//Arrange
	var tests = []struct {
		varName       string
		value         interface{}
		expectedValue string
	}{
		{"someString", "Zm9v", "foo"},
		{"someError", "1aaaaaZm9v", ""},
		{"someInt", 1, ""},
		{"someStruct", helper.TestStructure{TestString: "someMoreString"}, ""},
	}

	// Act on table
	for _, tt := range tests {
		testName := fmt.Sprintf("Getting String Return Correct Value -> %v,%v", tt.value, tt.varName)
		t.Run(testName, func(t *testing.T) {
			config := NewConfigService()
			config.UpsertKey(tt.varName, tt.value)

			// act
			keyValue := config.GetBase64(tt.varName)

			assert.Equal(t, tt.expectedValue, keyValue)
		})
	}
}

func TestGetBase64_WithNotFoundValue_ReturnEmptyString(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetBase64("notFound")

	assert.Equal(t, "", keyValue)
}

func TestGetBase64_WithEmptyKey_ReturnEmptyString(t *testing.T) {
	//Arrange
	config := NewConfigService()

	// act
	keyValue := config.GetBase64("")

	assert.Equal(t, "", keyValue)
}
