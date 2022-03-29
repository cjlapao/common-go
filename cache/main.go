package cache

import (
	"reflect"
)

var globalCacheService *CacheService

type CacheProvider interface {
	Get(name string) *interface{}
	Set(name string, value interface{})
}

type CacheService struct {
	Providers []CacheProvider
}

func New() *CacheService {
	globalCacheService = &CacheService{
		Providers: make([]CacheProvider, 0),
	}

	return globalCacheService
}

func Get() *CacheService {
	if globalCacheService != nil {
		return globalCacheService
	}

	return New()
}

func (c *CacheService) RegisterProvider(providers ...CacheProvider) {
	for _, registerProvider := range providers {
		found := false
		for _, provider := range c.Providers {
			if reflect.TypeOf(provider) == reflect.TypeOf(registerProvider) {
				found = true
				break
			}
		}
		if !found {
			c.Providers = append(c.Providers, registerProvider)
		}
	}
}
