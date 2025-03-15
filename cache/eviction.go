package cache

import "container/list"

type EvictionStrategy interface {
	Evict(cache *Cache)
}

type LRU struct{}

func (l LRU) Evict(cache *Cache) {
	if elem := cache.order.Front(); elem != nil {
		cache.Remove(elem)
	}
}

type LFU struct{}

func (l LFU) Evict(cache *Cache) {
	var leastUsed *list.Element
	for e := cache.order.Front(); e != nil; e = e.Next() {
		item := e.Value.(*CacheItem)
		if leastUsed == nil || item.frequency < leastUsed.Value.(*CacheItem).frequency {
			leastUsed = e
		}
	}
	if leastUsed != nil {
		cache.order.Remove(leastUsed)
		delete(cache.items, leastUsed.Value.(*CacheItem).key)
	}
}
