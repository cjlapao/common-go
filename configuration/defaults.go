package configuration

func (c *ConfigurationService) RegisterDefaults() *ConfigurationService {
	c.RegisterProvider(CachedVaultConfigurationProvider{}, EnvironmentConfigurationProvider{})

	return c
}
