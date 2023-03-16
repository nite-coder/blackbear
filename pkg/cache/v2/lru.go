package v2

import (
	"container/list"
	"sync"
)

// CacheItem is an item in the LRU cache
type CacheItem[V any] struct {
	key   string
	value V
}

// LRUCache is a generic type Least Recently Used (LRU) cache
type LRUCache[V any] struct {
	maxItems int
	items    map[string]*list.Element
	list     *list.List
	mutex    sync.Mutex
}

// NewLRUCache returns a new instance of the LRUCache with specified capacity
func NewLRUCache[V any](maxItems int) *LRUCache[V] {
	return &LRUCache[V]{
		maxItems: maxItems,
		items:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Len returns the number of items in the cache
func (c *LRUCache[V]) Len() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.list.Len()
}

// Get retrieves an item from the cache by key. Returns the value and true if the item exists, otherwise false
func (c *LRUCache[V]) Get(key string) (V, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*CacheItem[V]).value, true
	}

	var result V
	return result, false
}

// Put adds an item to the cache. If the item already exists, update its value and move it to the front of the list.
// If the cache is full, remove the least recently used item before adding the new item.
func (c *LRUCache[V]) Put(key string, value V) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*CacheItem[V]).value = value
		return
	}

	item := &CacheItem[V]{key: key, value: value}

	if c.list.Len() >= c.maxItems {
		// remove the least recently used item
		elem := c.list.Back()
		item := elem.Value.(*CacheItem[V])
		delete(c.items, item.key)
		c.list.Remove(elem)
	}

	elem := c.list.PushFront(item)
	c.items[key] = elem
}

// Delete removes an item from the cache by key. Returns true if the item exists and removed, otherwise false.
func (c *LRUCache[V]) Delete(key string) bool {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if elem, ok := c.items[key]; ok {
		c.list.Remove(elem)
		delete(c.items, key)
		return true
	}

	return false
}

// Clear removes all items from the cache
func (c *LRUCache[V]) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.list.Init()
	c.items = make(map[string]*list.Element)
}
