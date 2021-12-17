package jwt_token_cache

import (
	"strings"
)

var globalTokenCache *TokenCacheService

type TokenCacheItem struct {
	Name  string
	Token CachedJwtToken
}

type TokenCacheService struct {
	Items []TokenCacheItem
}

func New() *TokenCacheService {
	if globalTokenCache != nil {
		return globalTokenCache
	}

	globalTokenCache = &TokenCacheService{
		Items: make([]TokenCacheItem, 0),
	}

	return globalTokenCache
}

func (c *TokenCacheService) Get(name string) *CachedJwtToken {
	for _, item := range c.Items {
		if strings.EqualFold(name, item.Name) {
			return &item.Token
		}
	}

	return nil
}

func (c *TokenCacheService) Set(name string, token CachedJwtToken) {
	found := false
	for _, item := range c.Items {
		if strings.EqualFold(name, item.Name) {
			found = true
			break
		}
	}

	if !found {
		c.Items = append(c.Items, TokenCacheItem{
			Name:  name,
			Token: token,
		})
	}
}
