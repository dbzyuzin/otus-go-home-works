package hw04_lru_cache //nolint:golint,stylecheck

type Key string

type Cache interface {
	Set(key string, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key string) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                                 // Очистить кэш
}

type lruCache struct {
	cap      int                // - capacity
	queue    List               // - queue
	cacheMap map[Key]*cacheItem // - items
}

type cacheItem struct {
	key   string
	value interface{}
}

func (l lruCache) Set(key string, value interface{}) bool {
	panic("implement me")
}

func (l lruCache) Get(key string) (interface{}, bool) {
	panic("implement me")
}

func (l lruCache) Clear() {
	panic("implement me")
}

func NewCache(capacity int) Cache {
	return &lruCache{}
}
