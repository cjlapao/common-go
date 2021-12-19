package service_provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceProviderShouldReturnDefaultServicesWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	svc := New()

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
}

func TestNewServiceProviderShouldResetContainerAndReturnDefaultServicesWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	svc := New()

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
}

func TestGetServiceProviderShouldReturnNewServiceProviderWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	svc := Get()

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
}

func TestGetServiceProviderShouldReturnExistingServicesWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	svc := Get()

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
}
