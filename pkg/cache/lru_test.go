package cache

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLRU(t *testing.T) {
	lru := NewLRU(2)
	lru.Add("key1", "val1")
	lru.Add("key2", "val2")
	lru.Add("key3", "val3")

	val2, ok := lru.Get("key2")
	assert.Equal(t, "val2", val2)
	assert.Equal(t, true, ok)

	val1, ok := lru.Get("key1")
	assert.Equal(t, nil, val1)
	assert.Equal(t, false, ok)
}
