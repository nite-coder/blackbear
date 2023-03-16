package v2

import (
	"sync"
	"time"
)

type item[V any] struct {
	value      V
	expiration int
}

type Cache[K comparable, V any] struct {
	store       map[K]item[V]
	mutex       sync.RWMutex
	cleanupTime time.Duration
	stopCleanup chan bool
}

var NoExpiration time.Duration = 0

func NewCache[K comparable, V any](cleanupTime time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		store:       make(map[K]item[V]),
		cleanupTime: cleanupTime,
		stopCleanup: make(chan bool),
	}

	if cleanupTime > 0 {
		go cache.cleanupExpired()
	}

	return cache
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.store[key]
	if !ok {
		return item.value, false
	}

	if item.expiration != 0 && item.expiration < int(time.Now().Unix()) {
		delete(c.store, key)
		return item.value, false
	}

	return item.value, true
}

func (c *Cache[K, V]) Set(key K, value V, expiration time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	var exp int
	if expiration == NoExpiration {
		exp = 0
	} else {
		exp = int(time.Now().Add(expiration).Unix())
	}

	c.store[key] = item[V]{value, exp}
}

func (c *Cache[K, V]) Keys() []K {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	keys := make([]K, 0, len(c.store))
	for key := range c.store {
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache[K, V]) Delete(key K) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.store, key)
}

func (c *Cache[K, V]) Count() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return len(c.store)
}

func (c *Cache[K, V]) cleanupExpired() {
	ticker := time.NewTicker(c.cleanupTime)
	for {
		select {
		case <-ticker.C:
			c.mutex.Lock()
			for key, item := range c.store {
				if item.expiration != 0 && item.expiration < int(time.Now().Unix()) {
					delete(c.store, key)
				}
			}
			c.mutex.Unlock()
		case <-c.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

func (c *Cache[K, V]) StopCleanup() {
	c.stopCleanup <- true
}
