package cache

import "time"

type CacheItem struct {
	key       string
	value     any
	timestamp time.Time
	frequency int
}
