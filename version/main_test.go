package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromStringShortVersion(t *testing.T) {
	// arrange
	ver := "v0.2.4"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, err, "Error should be null")
	assert.Equalf(t, 0, v.Major, "Major should be 0")
	assert.Equalf(t, 2, v.Minor, "Minor should be 2")
	assert.Equalf(t, 4, v.Build, "Build should be 4")
	assert.Equalf(t, 0, v.Rev, "Rev should be 0")
}

func TestFromStringLongVersion(t *testing.T) {
	// arrange
	ver := "v1.2.4.7"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, err, "Error should be null")
	assert.Equalf(t, 1, v.Major, "Major should be 1")
	assert.Equalf(t, 2, v.Minor, "Minor should be 2")
	assert.Equalf(t, 4, v.Build, "Build should be 4")
	assert.Equalf(t, 7, v.Rev, "Rev should be 7")
}

func TestFromStringJustMajorVersion(t *testing.T) {
	// arrange
	ver := "1"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, err, "Error should be null")
	assert.Equalf(t, 1, v.Major, "Major should be 1")
	assert.Equalf(t, 0, v.Minor, "Minor should be 0")
	assert.Equalf(t, 0, v.Build, "Build should be 0")
	assert.Equalf(t, 0, v.Rev, "Rev should be 0")
}

func TestFromStringWithMajorAndMinorVersion(t *testing.T) {
	// arrange
	ver := "1.2"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, err, "Error should be null")
	assert.Equalf(t, 1, v.Major, "Major should be 1")
	assert.Equalf(t, 2, v.Minor, "Minor should be 2")
	assert.Equalf(t, 0, v.Build, "Build should be 0")
	assert.Equalf(t, 0, v.Rev, "Rev should be 0")
}

func TestFromStringWithMajorError(t *testing.T) {
	// arrange
	ver := "a.2"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, v, "Version should be null")
	assert.NotNilf(t, err, "Error should not be null")
	assert.ErrorContainsf(t, err, "could not parse major version", "Error message should be about major")
}

func TestFromStringWithMinorError(t *testing.T) {
	// arrange
	ver := "1.b"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, v, "Version should be null")
	assert.NotNilf(t, err, "Error should not be null")
	assert.ErrorContainsf(t, err, "could not parse minor version", "Error message should be about minor")
}

func TestFromStringWithBuildError(t *testing.T) {
	// arrange
	ver := "1.1.x"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, v, "Version should be null")
	assert.NotNilf(t, err, "Error should not be null")
	assert.ErrorContainsf(t, err, "could not parse build version", "Error message should be about build")
}

func TestFromStringWithRevError(t *testing.T) {
	// arrange
	ver := "1.1.1.x"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, v, "Version should be null")
	assert.NotNilf(t, err, "Error should not be null")
	assert.ErrorContainsf(t, err, "could not parse rev version", "Error message should be about rev")
}

func TestFromStringWithLongError(t *testing.T) {
	// arrange
	ver := "1.1.1.2.6"

	v, err := FromString(ver)

	//Assert
	assert.Nilf(t, v, "Version should be null")
	assert.NotNilf(t, err, "Error should not be null")
	assert.ErrorContainsf(t, err, "could not parse string", "Error message should be about parsing")
}

func TestToString(t *testing.T) {
	// arrange
	ver := "1.1.1.0"

	v, err := FromString(ver)

	//Assert
	assert.NotNilf(t, v, "Version should not be null")
	assert.Nilf(t, err, "Error should be nil")
	assert.Equalf(t, "1.1.1", v.String(), "ToString should be equal to 1.1.1")
}

func TestEmptyGet(t *testing.T) {
	// arrange
	v := Get()

	//Assert
	assert.NotNilf(t, v, "Version should not be null")
	assert.Equalf(t, 0, v.Major, "Major should be 0")
	assert.Equalf(t, 0, v.Minor, "Minor should be 0")
	assert.Equalf(t, 0, v.Build, "Build should be 0")
	assert.Equalf(t, 0, v.Rev, "Rev should be 0")
}

func TestGetMajor(t *testing.T) {
	// arrange
	appVersion = nil
	v := Get(1)

	//Assert
	assert.NotNilf(t, v, "Version should not be null")
	assert.Equalf(t, 1, v.Major, "Major should be 1")
	assert.Equalf(t, 0, v.Minor, "Minor should be 0")
	assert.Equalf(t, 0, v.Build, "Build should be 0")
	assert.Equalf(t, 0, v.Rev, "Rev should be 0")
}

func TestGetMajorMinor(t *testing.T) {
	// arrange
	appVersion = nil
	v := Get(1, 2)

	//Assert
	assert.NotNilf(t, v, "Version should not be null")
	assert.Equalf(t, 1, v.Major, "Major should be 1")
	assert.Equalf(t, 2, v.Minor, "Minor should be 2")
	assert.Equalf(t, 0, v.Build, "Build should be 0")
	assert.Equalf(t, 0, v.Rev, "Rev should be 0")
}

func TestGetMajorMinorBuild(t *testing.T) {
	// arrange
	appVersion = nil
	v := Get(1, 2, 3)

	//Assert
	assert.NotNilf(t, v, "Version should not be null")
	assert.Equalf(t, 1, v.Major, "Major should be 1")
	assert.Equalf(t, 2, v.Minor, "Minor should be 2")
	assert.Equalf(t, 3, v.Build, "Build should be 3")
	assert.Equalf(t, 0, v.Rev, "Rev should be 0")
}

func TestGetMajorMinorBuildRev(t *testing.T) {
	// arrange
	appVersion = nil
	v := Get(1, 2, 3, 4)

	//Assert
	assert.NotNilf(t, v, "Version should not be null")
	assert.Equalf(t, 1, v.Major, "Major should be 1")
	assert.Equalf(t, 2, v.Minor, "Minor should be 2")
	assert.Equalf(t, 3, v.Build, "Build should be 3")
	assert.Equalf(t, 4, v.Rev, "Rev should be 4")
}

func TestGetAlwaysReturnInstantiated(t *testing.T) {
	// arrange
	appVersion = nil
	v := Get(1, 2, 3, 4)
	v1 := Get(2, 3, 4, 5)

	//Assert
	assert.NotNilf(t, v, "Version should not be null")
	assert.Equalf(t, v.Major, v1.Major, "Major should be 1")
	assert.Equalf(t, v.Minor, v1.Minor, "Minor should be 2")
	assert.Equalf(t, v.Build, v1.Build, "Build should be 3")
	assert.Equalf(t, v.Rev, v1.Rev, "Rev should be 4")
}
