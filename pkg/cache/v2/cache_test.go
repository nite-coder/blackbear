package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache := NewCache[string, int](time.Minute)

	// Set a value
	cache.PutWithTTL("foo", 42, time.Second)

	// Check that the value is retrievable
	val, ok := cache.Get("foo")
	assert.True(t, ok)
	assert.Equal(t, 42, val)

	// Wait for the expiration of the value
	time.Sleep(2 * time.Second)

	// Check that the value is no longer retrievable
	_, ok = cache.Get("foo")
	assert.False(t, ok)

	// Set a value with no expiration
	cache.Put("bar", 99)

	// Check that the value is still retrievable after some time
	time.Sleep(time.Second)
	val, ok = cache.Get("bar")
	assert.True(t, ok)
	assert.Equal(t, 99, val)

	// Delete the value and check that it's gone
	cache.Delete("bar")
	_, ok = cache.Get("bar")
	assert.False(t, ok)

	// Check that the cache count is correct
	assert.Equal(t, 0, cache.Count())
	cache.PutWithTTL("baz", 123, time.Minute)
	assert.Equal(t, 1, cache.Count())

	// Stop the cache cleanup routine and check that it stops
	cache.StopCleanup()
	time.Sleep(time.Second)
	cache.PutWithTTL("qux", 456, time.Minute)
	assert.Equal(t, 2, cache.Count())
}
