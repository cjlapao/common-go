package service_provider

import (
	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/constants"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
)

type ServiceProvider struct {
	Configuration *configuration.ConfigurationService
	Version       *version.Version
	Logger        *log.Logger
}

var globalProviderContainer *ServiceProvider

func New() *ServiceProvider {
	if globalProviderContainer != nil {
		globalProviderContainer = nil
	}

	globalProviderContainer = &ServiceProvider{}
	globalProviderContainer.Logger = log.Get()
	globalProviderContainer.Version = version.Get()
	globalProviderContainer.Configuration = configuration.Get()
	globalProviderContainer.Logger.WithTimestamp()
	if globalProviderContainer.Configuration.GetBool(constants.DEBUG_ENVIRONMENT) {
		globalProviderContainer.Logger.WithDebug()
	}
	if globalProviderContainer.Configuration.GetBool(constants.TRACE_ENVIRONMENT) {
		globalProviderContainer.Logger.WithTrace()
	}

	return globalProviderContainer
}

func Get() *ServiceProvider {
	if globalProviderContainer != nil {
		return globalProviderContainer
	}

	return New()
}
