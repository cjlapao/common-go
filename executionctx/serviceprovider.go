package executionctx

import (
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
)

type ServiceProvider struct {
	Context *Context
	Version *version.Version
	Logger  *log.Logger
}

var globalProviderContainer *ServiceProvider

func NewServiceProvider() *ServiceProvider {
	if globalProviderContainer != nil {
		globalProviderContainer = nil
		NewContext()
	}

	globalProviderContainer = &ServiceProvider{}
	globalProviderContainer.Context = GetContext()
	globalProviderContainer.Logger = log.Get()
	globalProviderContainer.Version = version.Get()
	globalProviderContainer.Logger.UseTimestamp = true
	return globalProviderContainer
}

func GetServiceProvider() *ServiceProvider {
	if globalProviderContainer != nil {
		return globalProviderContainer
	}

	return NewServiceProvider()
}
