package cache

import (
	"runtime"
	"sync"
	"time"
)

const (
	// For use with functions that take an expiration time.
	NoExpiration time.Duration = -1
	// For use with functions that take an expiration time. Equivalent to
	// passing in the same expiration duration as was given to New() or
	// NewFrom() when the cache was created (e.g. 5 minutes.)
	DefaultExpiration time.Duration = 0
)

type item struct {
	object     interface{}
	expiration int64
}

type Cache struct {
	defaultExpiration time.Duration
	items             map[string]item
	mu                sync.RWMutex
	janitor           *janitor
}

// Return a new cache with a given default expiration duration and cleanup
// interval. If the expiration duration is less than one (or NoExpiration),
// the items in the cache never expire (by default), and must be deleted
// manually. If the cleanup interval is less than one, expired items are not
// deleted from the cache before calling c.DeleteExpired().
func New(defaultExpiration, cleanupInterval time.Duration) *Cache {
	c := Cache{
		defaultExpiration: defaultExpiration,
		items:             map[string]item{},
	}

	if cleanupInterval > 0 {
		runJanitor(&c, cleanupInterval)
		runtime.SetFinalizer(&c, stopJanitor)
	}

	return &c
}

// Add an item to the cache, replacing any existing item. If the duration is 0
// (DefaultExpiration), the cache's default expiration time is used. If it is -1
// (NoExpiration), the item never expires.
func (c *Cache) Set(key string, val interface{}, d time.Duration) {
	var e int64

	if d == DefaultExpiration {
		d = c.defaultExpiration
	}

	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[key] = item{
		object:     val,
		expiration: e,
	}
}

// Add an item to the cache, replacing any existing item, using the default
// expiration.
func (c *Cache) SetDefaultExpiration(k string, x interface{}) {
	c.Set(k, x, DefaultExpiration)
}

// Get an item from the cache. Returns the item or nil, and a bool indicating
// whether the key was found.
func (c *Cache) Get(k string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[k]

	if !found {
		return nil, false
	}

	if item.expiration > 0 {
		if time.Now().UnixNano() > item.expiration {
			return nil, false
		}
	}

	return item.object, true
}

// Returns the number of items in the cache. This may include items that have
// expired, but have not yet been cleaned up.
func (c *Cache) Count() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	n := len(c.items)
	return n
}

// Delete all expired items from the cache.
func (c *Cache) DeleteExpired() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if v.expiration > 0 && now > v.expiration {
			delete(c.items, k)
		}
	}
}

// DeleteAll removes all items from the cache.
func (c *Cache) DeleteAll() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items = map[string]item{}
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (j *janitor) Run(c *Cache) {
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			c.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return
		}
	}
}

func stopJanitor(c *Cache) {
	c.janitor.stop <- true
}

func runJanitor(c *Cache, ci time.Duration) {
	j := &janitor{
		Interval: ci,
		stop:     make(chan bool),
	}
	c.janitor = j
	go j.Run(c)
}
