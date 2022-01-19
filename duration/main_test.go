package duration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItNormalizesMS(t *testing.T) {
	// Arrange
	duration := &Duration{
		MilliSeconds: 1001,
	}

	// Act
	result := duration.String()

	// Assert
	assert.Equal(t, 1, duration.Seconds)
	assert.Equal(t, 1, duration.MilliSeconds)
	assert.Equal(t, "PT1S", result)
}

func TestItNormalizesS(t *testing.T) {
	// Arrange
	duration := &Duration{
		Seconds: 61,
	}

	// Act
	result := duration.String()

	// Assert
	assert.Equal(t, 1, duration.Seconds)
	assert.Equal(t, 1, duration.Minutes)
	assert.Equal(t, "PT1M1S", result)
}

func TestItNormalizesM(t *testing.T) {
	// Arrange
	duration := &Duration{
		Minutes: 61,
	}

	// Act
	result := duration.String()

	// Assert
	assert.Equal(t, 1, duration.Hours)
	assert.Equal(t, 1, duration.Minutes)
	assert.Equal(t, "PT1H1M", result)
}

func TestItNormalizesH(t *testing.T) {
	// Arrange
	duration := &Duration{
		Hours: 25,
	}

	// Act
	result := duration.String()

	// Assert
	assert.Equal(t, 1, duration.Hours)
	assert.Equal(t, 1, duration.Days)
	assert.Equal(t, "P1DT1H", result)
}

func TestItNormalizesD(t *testing.T) {
	// Arrange
	duration := &Duration{
		Days: 8,
	}

	// Act
	result := duration.String()

	// Assert
	assert.Equal(t, 1, duration.Weeks)
	assert.Equal(t, 1, duration.Days)
	assert.Equal(t, "P1W1D", result)
}

func TestItDoesntNormalizesW(t *testing.T) {
	// Arrange
	duration := &Duration{
		Days: 56,
	}

	// Act
	result := duration.String()

	// Assert
	assert.Equal(t, 56, duration.Weeks)
	assert.Equal(t, 0, duration.Years)
	assert.Equal(t, "P56W", result)
}
