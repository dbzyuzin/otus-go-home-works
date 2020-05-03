package hw04_lru_cache //nolint:golint,stylecheck

import "sync"

type Key string

type Cache interface {
	Set(key string, value interface{}) bool // Добавить значение в кэш по ключу
	Get(key string) (interface{}, bool)     // Получить значение из кэша по ключу
	Clear()                                 // Очистить кэш
}

type lruCache struct {
	mux      sync.RWMutex         // - mutex
	cap      int                  // - capacity
	queue    List                 // - queue
	cacheMap map[string]*listItem // - items
}

type cacheItem struct {
	key   string
	value interface{}
}

func (l *lruCache) Set(key string, value interface{}) bool {
	l.mux.Lock()
	defer l.mux.Unlock()
	lstItem, ok := l.cacheMap[key]
	if ok {
		lstItem.Value = cacheItem{
			key:   key,
			value: value,
		}
		l.queue.MoveToFront(lstItem)
	} else {
		lstItem = l.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
		l.cacheMap[key] = lstItem

		if l.queue.Len() > l.cap {
			back := l.queue.Back()
			l.queue.Remove(back)
			delete(l.cacheMap, back.Value.(cacheItem).key)
		}
	}
	return ok
}

func (l *lruCache) Get(key string) (interface{}, bool) {
	l.mux.RLock()
	lstItem, ok := l.cacheMap[key]
	l.mux.RUnlock()
	if ok {
		l.mux.Lock()
		defer l.mux.Unlock()
		l.queue.MoveToFront(lstItem)
		cache := lstItem.Value.(cacheItem)
		return cache.value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.queue = NewList()
	l.cacheMap = make(map[string]*listItem, l.cap)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		cap:      capacity,
		queue:    NewList(),
		cacheMap: make(map[string]*listItem, capacity),
	}
}
