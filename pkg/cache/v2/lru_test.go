package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRUCache(t *testing.T) {
	// Create a new LRU cache with a capacity of 2 items
	cache := NewLRUCache[int](2)

	// Add some items to the cache
	cache.Put("key1", 1)
	cache.Put("key2", 2)

	// Check the length of the cache
	assert.Equal(t, 2, cache.Len())

	// Retrieve an item from the cache
	value, ok := cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	// Add another item to the cache
	cache.Put("key3", 3)

	// Check that the least recently used item was removed
	_, ok = cache.Get("key2")
	assert.False(t, ok)

	// Check that the remaining items are still in the cache
	value, ok = cache.Get("key1")
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	value, ok = cache.Get("key3")
	assert.True(t, ok)
	assert.Equal(t, 3, value)

	// Remove an item from the cache
	ok = cache.Delete("key1")
	assert.True(t, ok)

	// Check that the removed item is no longer in the cache
	_, ok = cache.Get("key1")
	assert.False(t, ok)

	// Clear the cache
	cache.Clear()

	// Check that the cache is empty
	assert.Equal(t, 0, cache.Len())
}
