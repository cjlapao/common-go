package executionctx

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServiceProviderShouldReturnDefaultServicesWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	svc := NewServiceProvider()

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Context, "Context Should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
}

func TestNewServiceProviderShouldResetContainerAndReturnDefaultServicesWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	oldSvc := NewServiceProvider()
	oldContextId := oldSvc.Context.CorrelationId
	svc := NewServiceProvider()
	currentContextId := svc.Context.CorrelationId

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Context, "Context Should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
	assert.NotEqualf(t, oldContextId, currentContextId, "Correlation Id should not be the same")
}

func TestGetServiceProviderShouldReturnNewServiceProviderWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	svc := GetServiceProvider()

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Context, "Context Should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
}

func TestGetServiceProviderShouldReturnExistingServicesWithDefaultValues(t *testing.T) {
	// arrange
	globalProviderContainer = nil
	oldSvc := NewServiceProvider()
	oldContextId := oldSvc.Context.CorrelationId
	svc := GetServiceProvider()
	currentContextId := svc.Context.CorrelationId

	//Assert
	assert.NotNilf(t, svc, "Service Provider should not be nil")
	assert.NotNilf(t, svc.Logger, "Logger should not be nil")
	assert.NotNilf(t, svc.Context, "Context Should not be nil")
	assert.NotNilf(t, svc.Version, "Version Service should not be nil")
	assert.Equalf(t, oldContextId, currentContextId, "Correlation Id should not be the same")
}
