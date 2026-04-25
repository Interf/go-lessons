package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLRUCache(t *testing.T) {

	tests := map[string]struct {
		capacity int
		want     int
	}{
		"capacity 2": {
			capacity: 2,
			want:     2,
		},
		"capacity 0": {
			capacity: 0,
			want:     8,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			cache := NewLRUCache(tt.capacity)

			assert.Equal(t, tt.want, cache.capacity)
		})
	}
}

func TestLRUCacheGet(t *testing.T) {
	cache := NewLRUCache(3)

	for i := 0; i < 3; i++ {
		cache.Put(i, i)
	}

	tests := map[string]struct {
		key  int
		want int
		ok   bool
	}{
		"exist key": {
			key:  2,
			want: 2,
			ok:   true,
		},
		"not exist key": {
			key:  55,
			want: 0,
			ok:   false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, ok := cache.Get(tt.key)

			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.ok, ok)
		})
	}
}

func TestLRUCachePut(t *testing.T) {
	cache := NewLRUCache(1)

	tests := []struct {
		name    string
		key     int
		value   int
		findKey int
		want    int
	}{
		{
			name:    "new element",
			key:     1,
			value:   1,
			findKey: 1,
			want:    1,
		},
		{
			name:    "update element",
			key:     1,
			value:   5,
			findKey: 1,
			want:    5,
		},
		{
			name:    "delete tail",
			key:     55,
			value:   55,
			findKey: 1,
			want:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cache.Put(tt.key, tt.value)
			got, _ := cache.Get(tt.findKey)

			assert.Equal(t, tt.want, got)
		})
	}

}
