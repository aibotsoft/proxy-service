package cache

import (
	"github.com/dgraph-io/ristretto"
)

func NewCache() (*ristretto.Cache, error) {
	return ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e6,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Metrics:     true,
	})
}
