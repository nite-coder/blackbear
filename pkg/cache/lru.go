package cache

import (
	"container/list"
	"sync"
)

// LRUCache is a LRU cache.
type LRUCache struct {
	maxItems int64
	count    int64
	ll       *list.List
	mu       sync.RWMutex
	cache    map[string]*list.Element
}

// LRUCacher is a LRU cache interface
type LRUCacher interface {
	Put(key string, val interface{})
	Get(key string) (interface{}, bool)
}

type entry struct {
	key   string
	value interface{}
}

// New is the Constructor of Cache
func NewLRU(maxItems int64) *LRUCache {
	return &LRUCache{
		maxItems: maxItems,
		ll:       list.New(),
		mu:       sync.RWMutex{},
		cache:    make(map[string]*list.Element),
	}
}

// Put adds a value to the cache.
func (c *LRUCache) Put(key string, val interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if ele, ok := c.cache[key]; ok {
		ele.Value = &entry{key: key, value: val}
		c.ll.MoveToFront(ele)
	} else {
		kv := &entry{key: key, value: val}
		ele := c.ll.PushFront(kv)
		c.cache[key] = ele
		c.count++
	}

	if c.count > c.maxItems {
		c.removeOldest()
	}
}

// Get returns a value from cache
func (c *LRUCache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (c *LRUCache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
	}
}
