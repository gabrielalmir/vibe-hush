package cache

import (
	"container/list"
	"time"
)

type SetMessage struct {
	key   string
	value any
}

type GetMessage struct {
	key      string
	response chan any
}

type DeleteMessage struct {
	key string
}

type GetAllMessage struct {
	response chan map[string]any
}

type Cache struct {
	capacity   int
	expiration time.Duration
	strategy   EvictionStrategy
	items      map[string]*list.Element
	order      *list.List
	commands   chan any
}

func NewCache(capacity int, expiration time.Duration, strategy EvictionStrategy) *Cache {
	c := &Cache{
		capacity:   capacity,
		expiration: expiration,
		strategy:   strategy,
		items:      make(map[string]*list.Element),
		order:      list.New(),
		commands:   make(chan any),
	}
	go c.run()
	return c
}

func (c *Cache) run() {
	for msg := range c.commands {
		switch m := msg.(type) {
		case SetMessage:
			c.handleSet(m.key, m.value)
		case GetMessage:
			m.response <- c.handleGet(m.key)
		case GetAllMessage:
			m.response <- c.handleGetAll()
		case DeleteMessage:
			c.handleDelete(m.key)
		}
	}
}

func (c *Cache) Set(key string, value any) {
	c.commands <- SetMessage{key: key, value: value}
}

func (c *Cache) Get(key string) any {
	resp := make(chan any)
	c.commands <- GetMessage{key: key, response: resp}
	return <-resp
}

func (c *Cache) GetAll() map[string]any {
	resp := make(chan map[string]any)
	c.commands <- GetAllMessage{response: resp}
	return <-resp
}

func (c *Cache) Delete(key string) {
	c.commands <- DeleteMessage{key: key}
}

func (c *Cache) handleSet(key string, value any) {
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

func (c *Cache) handleGet(key string) any {
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

func (c *Cache) handleGetAll() map[string]any {
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

func (c *Cache) handleDelete(key string) {
	if elem, found := c.items[key]; found {
		c.order.Remove(elem)
		delete(c.items, key)
	}
}

func (c *Cache) Remove(elem *list.Element) {
	delete(c.items, elem.Value.(*CacheItem).key)
	c.order.Remove(elem)
}
