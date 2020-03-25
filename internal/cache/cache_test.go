package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache, err := NewCache()
	assert.Nil(t, err)
	cache.Set("key", "value", 0)
	time.Sleep(1 * time.Millisecond)

	value, ok := cache.Get("key")
	assert.Equal(t, true, ok)
	assert.Equal(t, "value", value)
}
