package cache

import (
	"container/list"
	"sync"
	"time"
)

type Cache struct {
	capacity   int
	expiration time.Duration
	strategy   EvictionStrategy
	items      map[string]*list.Element
	order      *list.List
	mutex      sync.Mutex
}

func NewCache(capacity int, expiration time.Duration, strategy EvictionStrategy) *Cache {
	return &Cache{
		capacity:   capacity,
		expiration: expiration,
		strategy:   strategy,
		items:      make(map[string]*list.Element),
		order:      list.New(),
	}
}

func (c *Cache) Get(key string) any {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, found := c.items[key]; found {
		item := elem.Value.(*CacheItem)
		if time.Since(item.timestamp) > c.expiration {
			c.order.Remove(elem)
			delete(c.items, key)
			return nil
		}
		if _, ok := c.strategy.(LRU); ok {
			c.order.MoveToBack(elem)
		} else if _, ok := c.strategy.(LFU); ok {
			item.frequency++
		}
		return item.value
	}
	return nil
}

func (c *Cache) GetAll() map[string]any {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	result := make(map[string]any)
	for key, elem := range c.items {
		item := elem.Value.(*CacheItem)
		if time.Since(item.timestamp) <= c.expiration {
			result[key] = item.value
		} else {
			c.order.Remove(elem)
			delete(c.items, key)
		}
	}
	return result
}

func (c *Cache) Set(key string, value any) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, found := c.items[key]; found {
		c.order.MoveToBack(elem)
		elem.Value.(*CacheItem).value = value
		elem.Value.(*CacheItem).timestamp = time.Now()
		return
	}

	if len(c.items) >= c.capacity {
		c.strategy.Evict(c)
	}

	item := &CacheItem{key, value, time.Now(), 1}
	elem := c.order.PushBack(item)
	c.items[key] = elem
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, found := c.items[key]; found {
		c.order.Remove(elem)
		delete(c.items, key)
	}
}
