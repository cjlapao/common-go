package cache

type CacheService interface {
	Get(name string) interface{}
	Set(name string, value interface{})
}
