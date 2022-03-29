package execution_context

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContextShouldReturnContextWithDefaultValues(t *testing.T) {
	// arrange
	contextService = nil
	ctx, err := New()

	//Assert
	assert.Nilf(t, err, "Error should be null")
	assert.NotNilf(t, ctx, "Context should not be nil")
	assert.NotNilf(t, ctx.Configuration, "Configuration should not be nil")
	assert.NotEmptyf(t, ctx.CorrelationId, "Correlation ID should not be empty")
	assert.Falsef(t, ctx.IsDevelopment, "IsDevelopment should be false")
	assert.Falsef(t, ctx.Debug, "Debug should be false")
	assert.Equalf(t, "Production", ctx.Environment, "Environment should be Production")
}

func TestInitNewContextShouldRunInit(t *testing.T) {
	// Arrange + Act
	contextService = nil
	ctx, err := InitNewContext(contextInitiation)

	// Assert
	assert.Nilf(t, err, "Error should be null")
	assert.NotNilf(t, ctx, "Context should not be nil")
	assert.NotNilf(t, ctx.Configuration, "Configuration should not be nil")
	assert.NotEmptyf(t, ctx.CorrelationId, "Correlation ID should not be empty")
	assert.Falsef(t, ctx.IsDevelopment, "IsDevelopment should be false")
	assert.Falsef(t, ctx.Debug, "Debug should be false")
	assert.Equalf(t, "Production", ctx.Environment, "Environment should be Production")
	assert.Equalf(t, "bar", ctx.Configuration.Get("foo"), "Initialization Key should be \"bar\"")
}

func TestInitNewContextShouldRunInitAndReturnError(t *testing.T) {
	// Arrange + Act
	contextService = nil
	ctx, err := InitNewContext(contextInitiationError)

	// Assert
	assert.NotNilf(t, err, "Error should be null")
	assert.EqualErrorf(t, err, "SomeRandomError", "Error message mismatch")
	assert.Nilf(t, ctx, "Context should not be nil")
}

func TestNewContextShouldResetOldContext(t *testing.T) {
	// Arrange + Act
	contextService = nil
	ctx, _ := New()
	ctx.IsDevelopment = true

	ctxNew, _ := New()

	// Assert
	assert.NotNil(t, ctxNew, "Context should not be nil")
	assert.Falsef(t, ctxNew.IsDevelopment, "Development should be false/default")
	assert.Falsef(t, ctxNew.Debug, "Debug should be false/default")
	assert.Equalf(t, "Production", ctx.Environment, "Environment should be \"Production\" found %v", ctx.Environment)
}

func TestNewContextShouldSetEnvironmentToDevelopment(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENVIRONMENT", "Development")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Truef(t, ctx.IsDevelopment, "Development should be true")
	assert.Equalf(t, "Development", ctx.Environment, "Environment should be \"Development\" found %v", ctx.Environment)
	os.Setenv("CJ_ENVIRONMENT", "")
}

func TestNewContextShouldSetEnvironmentToDevProd(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENVIRONMENT", "DevProd")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Falsef(t, ctx.IsDevelopment, "Development should be false")
	assert.Equalf(t, "DevProd", ctx.Environment, "Environment should be \"DevProd\" found %v", ctx.Environment)
	os.Setenv("CJ_ENVIRONMENT", "")
}

func TestNewContextShouldSetEnvironmentToCi(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENVIRONMENT", "Ci")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Truef(t, ctx.IsDevelopment, "Development should be false")
	assert.Equalf(t, "CI", ctx.Environment, "Environment should be \"CI\" found %v", ctx.Environment)
	os.Setenv("CJ_ENVIRONMENT", "")
}

func TestNewContextShouldSetEnvironmentToRelease(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENVIRONMENT", "Release")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Falsef(t, ctx.IsDevelopment, "Development should be false")
	assert.Equalf(t, "Release", ctx.Environment, "Environment should be \"Release\" found %v", ctx.Environment)
	os.Setenv("CJ_ENVIRONMENT", "")
}

func TestNewContextShouldSetEnvironmentToProduction(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENVIRONMENT", "Production")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Falsef(t, ctx.IsDevelopment, "Development should be false")
	assert.Equalf(t, "Production", ctx.Environment, "Environment should be \"Production\" found %v", ctx.Environment)
	os.Setenv("CJ_ENVIRONMENT", "")
}

func TestNewContextShouldSetEnvironmentToOthers(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENVIRONMENT", "RandomEnvironment")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Falsef(t, ctx.IsDevelopment, "Development should be false")
	assert.Equalf(t, "Production", ctx.Environment, "Environment should be \"Production\" found %v", ctx.Environment)
	os.Setenv("CJ_ENVIRONMENT", "")
}

func TestNewContextShouldSetDebugOn(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENABLE_DEBUG", "true")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Truef(t, ctx.Debug, "Debug should be true")
	os.Setenv("CJ_ENABLE_DEBUG", "")
}

func TestNewContextShouldSetDebugOff(t *testing.T) {
	// Arrange + Act
	contextService = nil
	os.Setenv("CJ_ENABLE_DEBUG", "enabled")
	ctx, _ := New()

	// Assert
	assert.NotNil(t, ctx, "Context should not be nil")
	assert.Falsef(t, ctx.Debug, "Debug should be false")
	os.Setenv("CJ_ENABLE_DEBUG", "")
}

func TestGetContextShouldReturnContextWithDefaultValues(t *testing.T) {
	// arrange
	contextService = nil
	ctx := Get()

	//Assert
	assert.NotNilf(t, ctx, "Context should not be nil")
	assert.NotNilf(t, ctx.Configuration, "Configuration should not be nil")
	assert.NotEmptyf(t, ctx.CorrelationId, "Correlation ID should not be empty")
	assert.Falsef(t, ctx.IsDevelopment, "IsDevelopment should be false")
	assert.Falsef(t, ctx.Debug, "Debug should be false")
	assert.Equalf(t, "Production", ctx.Environment, "Environment should be Production")
}

func TestGetContextShouldReturnExistingContext(t *testing.T) {
	// arrange
	existingCtx, _ := New()
	existingCtx.IsDevelopment = true
	existingCorrelationId := existingCtx.CorrelationId
	ctx := Get()

	//Assert
	assert.NotNilf(t, ctx, "Context should not be be nil")
	assert.NotNilf(t, ctx.Configuration, "Configuration should not be nil")
	assert.NotEmptyf(t, ctx.CorrelationId, "Correlation ID should not be empty")
	assert.Equalf(t, existingCorrelationId, ctx.CorrelationId, "Correlation ID should match old correlation id")
	assert.Falsef(t, ctx.Debug, "Debug should be false")
	assert.Equalf(t, "Production", ctx.Environment, "Environment should be Production")
	assert.Truef(t, ctx.IsDevelopment, "Development should be true")
}

func contextInitiation() error {
	os.Setenv("foo", "bar")

	return nil
}

func contextInitiationError() error {
	os.Setenv("foo", "bar")

	return errors.New("SomeRandomError")
}
