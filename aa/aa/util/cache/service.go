package cache

type CacheService interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
	Delete(string)
}
