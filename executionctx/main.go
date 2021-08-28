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

var globalProvider *ServiceProvider

func CreateProvider() *ServiceProvider {
	if globalProvider != nil {
		return globalProvider
	}

	globalProvider = &ServiceProvider{}
	globalProvider.Context = Get()
	globalProvider.Logger = log.Get()
	globalProvider.Version = version.Get()

	return globalProvider
}

func GetServiceProvider() *ServiceProvider {
	return globalProvider
}
