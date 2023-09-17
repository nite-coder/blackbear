package cache

import (
	"sync"
	"time"
)

var NoExpiration time.Duration = 0

type item[V any] struct {
	value      V
	expiration int
}

type Cache[K comparable, V any] struct {
	store       sync.Map
	cleanupTime time.Duration
	stopCleanup chan bool
}

// NewCache creates a new instance of Cache with the specified cleanup time.
func NewCache[K comparable, V any](cleanupTime time.Duration) *Cache[K, V] {
	cache := &Cache[K, V]{
		store:       sync.Map{},
		cleanupTime: cleanupTime,
		stopCleanup: make(chan bool),
	}

	if cleanupTime > 0 {
		go cache.cleanupExpired()
	}

	return cache
}

// Get retrieves a value from the cache for the specified key.
// Returns the value and a boolean indicating if the value was found.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	value, ok := c.store.Load(key)
	if !ok {
		var val V
		return val, false
	}

	item := value.(item[V])

	if item.expiration != 0 && item.expiration < int(time.Now().Unix()) {
		c.store.Delete(key)
		return item.value, false
	}

	return item.value, true
}

// PutWithTTL adds a value to the cache with the specified key and expiration time.
func (c *Cache[K, V]) PutWithTTL(key K, value V, ttl time.Duration) {
	var exp int
	if ttl == NoExpiration {
		exp = 0
	} else {
		exp = int(time.Now().Add(ttl).Unix())
	}

	c.store.Store(key, item[V]{value, exp})
}

// Put adds a value to the cache with the specified key and no expiration time.
func (c *Cache[K, V]) Put(key K, value V) {
	c.store.Store(key, item[V]{value, 0})
}

// Delete removes a value from the cache for the specified key.
func (c *Cache[K, V]) Delete(key K) {
	c.store.Delete(key)
}

// Keys returns a slice of all the keys in the cache.
func (c *Cache[K, V]) Keys() []K {
	keys := make([]K, 0)
	c.store.Range(func(key, value interface{}) bool {
		keys = append(keys, key.(K))
		return true
	})
	return keys
}

// Count returns the number of items currently stored in the cache.
func (c *Cache[K, V]) Count() int {
	count := 0
	c.store.Range(func(_, _ interface{}) bool {
		count++
		return true
	})
	return count
}

// StopCleanup stops the automatic cleanup of expired items from the cache.
func (c *Cache[K, V]) StopCleanup() {
	c.stopCleanup <- true
}

// cleanupExpired periodically checks for and removes expired items from the cache.
func (c *Cache[K, V]) cleanupExpired() {
	ticker := time.NewTicker(c.cleanupTime)
	for {
		select {
		case <-ticker.C:
			c.store.Range(func(key, value interface{}) bool {
				item := value.(item[V])
				if item.expiration != 0 && item.expiration < int(time.Now().Unix()) {
					c.store.Delete(key)
				}
				return true
			})
		case <-c.stopCleanup:
			ticker.Stop()
			return
		}
	}
}
